package service

import (
    "context"
    "encoding/json"
    "fmt"
    "sort"
    "strings"
    "sync"
    "time"

    appsv1 "k8s.io/api/apps/v1"
    corev1 "k8s.io/api/core/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

    "kubewatch-ai/internal/adapter/websocket"
    "kubewatch-ai/internal/core/model"
    "kubewatch-ai/internal/infrastructure/k8s"
    "kubewatch-ai/internal/infrastructure/monitoring"
)

type IncidentService struct {
    client          *k8s.Client
    metrics         *monitoring.PrometheusMetrics
    hub             *websocket.Hub
    mu              sync.RWMutex
    incidents       []model.Incident
    clusterOverview model.ClusterOverview
}

func NewIncidentService(client *k8s.Client, metrics *monitoring.PrometheusMetrics, hub *websocket.Hub) *IncidentService {
    return &IncidentService{client: client, metrics: metrics, hub: hub}
}

func (s *IncidentService) StartPolling(ctx context.Context) {
    ticker := time.NewTicker(15 * time.Second)
    defer ticker.Stop()

    s.refreshCluster(ctx)
    for {
        select {
        case <-ctx.Done():
            return
        case <-ticker.C:
            s.refreshCluster(ctx)
        }
    }
}

func (s *IncidentService) GetClusterOverview(ctx context.Context) (model.ClusterOverview, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    return s.clusterOverview, nil
}

func (s *IncidentService) GetNamespaces(ctx context.Context) ([]string, error) {
    namespaces, err := s.client.Clientset.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
    if err != nil {
        return nil, err
    }
    result := make([]string, 0, len(namespaces.Items))
    for _, item := range namespaces.Items {
        result = append(result, item.Name)
    }
    sort.Strings(result)
    return result, nil
}

func (s *IncidentService) GetIncidents(ctx context.Context) ([]model.Incident, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    out := make([]model.Incident, len(s.incidents))
    copy(out, s.incidents)
    return out, nil
}

func (s *IncidentService) refreshCluster(ctx context.Context) {
    incidents := make([]model.Incident, 0)
    incidentCounters := map[string]int{}
    namespaceSet := map[string]struct{}{}
    now := time.Now()

    pods, _ := s.client.Clientset.CoreV1().Pods("").List(ctx, metav1.ListOptions{})
    for _, pod := range pods.Items {
        if incident := s.buildPodIncident(pod.Namespace, pod.Name, pod.Status.ContainerStatuses); incident != nil {
            incidents = append(incidents, *incident)
            incidentCounters[incident.Type]++
            namespaceSet[incident.Namespace] = struct{}{}
        }
    }

    deployments, _ := s.client.Clientset.AppsV1().Deployments("").List(ctx, metav1.ListOptions{})
    for _, deploy := range deployments.Items {
        if incident := s.buildDeploymentIncident(deploy); incident != nil {
            incidents = append(incidents, *incident)
            incidentCounters[incident.Type]++
            namespaceSet[incident.Namespace] = struct{}{}
        }
    }

    services, _ := s.client.Clientset.CoreV1().Services("").List(ctx, metav1.ListOptions{})
    for _, svc := range services.Items {
        if incident := s.buildServiceIncident(ctx, svc); incident != nil {
            incidents = append(incidents, *incident)
            incidentCounters[incident.Type]++
            namespaceSet[incident.Namespace] = struct{}{}
        }
    }

    nodes, _ := s.client.Clientset.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
    for _, node := range nodes.Items {
        if incident := s.buildNodeIncident(node); incident != nil {
            incidents = append(incidents, *incident)
            incidentCounters[incident.Type]++
            namespaceSet[incident.Namespace] = struct{}{}
        }
    }

    overview := model.ClusterOverview{
        TotalNodes:          len(nodes.Items),
        ReadyNodes:          countReadyNodes(nodes.Items),
        TotalNamespaces:     len(namespaceSet),
        CrashLoopBackOff:    incidentCounters["CrashLoopBackOff"],
        UnhealthyDeployments: incidentCounters["UnhealthyDeployment"],
        ReplicaMismatch:     incidentCounters["DeploymentMismatch"],
        FailedServices:      incidentCounters["FailedService"] + incidentCounters["NodePressure"],
        ActiveIncidents:     len(incidents),
    }

    s.mu.Lock()
    s.incidents = incidents
    s.clusterOverview = overview
    s.mu.Unlock()

    s.metrics.UpdateIncidentCounts(len(incidents))
    s.metrics.UpdateOverviewMetrics(len(nodes.Items), len(pods.Items), len(deployments.Items), len(services.Items))

    payload, _ := json.Marshal(wsPayload{Incidents: incidents, Timestamp: now})
    s.hub.Broadcast(payload)
}

type wsPayload struct {
    Incidents []model.Incident `json:"incidents"`
    Timestamp time.Time        `json:"timestamp"`
}

func countReadyNodes(nodes []corev1.Node) int {
    ready := 0
    for _, node := range nodes {
        for _, condition := range node.Status.Conditions {
            if condition.Type == corev1.NodeReady && condition.Status == corev1.ConditionTrue {
                ready++
                break
            }
        }
    }
    return ready
}

func (s *IncidentService) buildPodIncident(namespace, name string, statuses []corev1.ContainerStatus) *model.Incident {
    for _, status := range statuses {
        if status.State.Waiting != nil && strings.Contains(status.State.Waiting.Reason, "CrashLoopBackOff") {
            return &model.Incident{
                ID:                  fmt.Sprintf("pod-%s-%s", namespace, name),
                Namespace:           namespace,
                Workload:            name,
                Type:                "CrashLoopBackOff",
                Summary:             "Pod is restarting repeatedly",
                Details:             fmt.Sprintf("%s: %s", status.Name, status.State.Waiting.Message),
                Severity:            model.SeverityHigh,
                SeverityScore:       78,
                SuggestedRemediation: "Inspect container logs and restart the failing pod or fix the startup command.",
                FirstSeen:           time.Now().Add(-5 * time.Minute),
                LastSeen:            time.Now(),
                Source:              "pod-monitor",
            }
        }
    }
    return nil
}

func (s *IncidentService) buildDeploymentIncident(deploy appsv1.Deployment) *model.Incident {
    if deploy.Status.Replicas > deploy.Status.ReadyReplicas {
        return &model.Incident{
            ID:                  fmt.Sprintf("deploy-%s-%s", deploy.Namespace, deploy.Name),
            Namespace:           deploy.Namespace,
            Workload:            deploy.Name,
            Type:                "DeploymentMismatch",
            Summary:             "Replica mismatch detected",
            Details:             fmt.Sprintf("desired=%d ready=%d", deploy.Status.Replicas, deploy.Status.ReadyReplicas),
            Severity:            model.SeverityMedium,
            SeverityScore:       55,
            SuggestedRemediation: "Review pod readiness and scaling policy, then reconcile deployment replicas.",
            FirstSeen:           time.Now().Add(-10 * time.Minute),
            LastSeen:            time.Now(),
            Source:              "deployment-monitor",
        }
    }
    return nil
}

func (s *IncidentService) buildServiceIncident(ctx context.Context, svc corev1.Service) *model.Incident {
    if svc.Spec.Selector == nil || len(svc.Spec.Selector) == 0 {
        return nil
    }
    endpoints, err := s.client.Clientset.CoreV1().Endpoints(svc.Namespace).Get(ctx, svc.Name, metav1.GetOptions{})
    if err != nil || len(endpoints.Subsets) == 0 {
        return &model.Incident{
            ID:                  fmt.Sprintf("service-%s-%s", svc.Namespace, svc.Name),
            Namespace:           svc.Namespace,
            Workload:            svc.Name,
            Type:                "FailedService",
            Summary:             "Service has no active endpoints",
            Details:             "The selected endpoints are not ready or missing.",
            Severity:            model.SeverityMedium,
            SeverityScore:       50,
            SuggestedRemediation: "Validate service selector labels and deployment readiness to restore traffic.",
            FirstSeen:           time.Now().Add(-15 * time.Minute),
            LastSeen:            time.Now(),
            Source:              "service-monitor",
        }
    }
    return nil
}

func (s *IncidentService) buildNodeIncident(node corev1.Node) *model.Incident {
    for _, condition := range node.Status.Conditions {
        if (condition.Type == corev1.NodeMemoryPressure || condition.Type == corev1.NodeDiskPressure || condition.Type == corev1.NodePIDPressure) && condition.Status == corev1.ConditionTrue {
            return &model.Incident{
                ID:                  fmt.Sprintf("node-%s-%s", node.Name, condition.Type),
                Namespace:           "cluster",
                Workload:            node.Name,
                Type:                "NodePressure",
                Summary:             "Node pressure condition detected",
                Details:             fmt.Sprintf("%s - %s", condition.Type, condition.Message),
                Severity:            model.SeverityHigh,
                SeverityScore:       72,
                SuggestedRemediation: "Evaluate node resource pressure and consider scaling or cordoning the node.",
                FirstSeen:           time.Now().Add(-20 * time.Minute),
                LastSeen:            time.Now(),
                Source:              "node-monitor",
            }
        }
    }
    return nil
}
