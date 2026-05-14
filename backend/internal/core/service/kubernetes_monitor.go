package service

import (
    "context"
    "fmt"
    "strings"
    "time"

    appsv1 "k8s.io/api/apps/v1"
    corev1 "k8s.io/api/core/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

    "kubewatch-ai/internal/core/model"
    "kubewatch-ai/internal/infrastructure/k8s"
)

type KubernetesMonitor struct {
    client *k8s.Client
}

func NewKubernetesMonitor(client *k8s.Client) *KubernetesMonitor {
    return &KubernetesMonitor{client: client}
}

func (m *KubernetesMonitor) ListNamespaces(ctx context.Context) ([]string, error) {
    namespaceList, err := m.client.Clientset.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
    if err != nil {
        return nil, err
    }

    namespaces := make([]string, 0, len(namespaceList.Items))
    for _, ns := range namespaceList.Items {
        namespaces = append(namespaces, ns.Name)
    }
    return namespaces, nil
}

func (m *KubernetesMonitor) ListPods(ctx context.Context, namespace string) ([]corev1.Pod, error) {
    pods, err := m.client.Clientset.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
    if err != nil {
        return nil, err
    }
    return pods.Items, nil
}

func (m *KubernetesMonitor) ListDeployments(ctx context.Context, namespace string) ([]appsv1.Deployment, error) {
    deployments, err := m.client.Clientset.AppsV1().Deployments(namespace).List(ctx, metav1.ListOptions{})
    if err != nil {
        return nil, err
    }
    return deployments.Items, nil
}

func (m *KubernetesMonitor) ListNodes(ctx context.Context) ([]corev1.Node, error) {
    nodes, err := m.client.Clientset.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
    if err != nil {
        return nil, err
    }
    return nodes.Items, nil
}

func (m *KubernetesMonitor) ListServices(ctx context.Context, namespace string) ([]corev1.Service, error) {
    services, err := m.client.Clientset.CoreV1().Services(namespace).List(ctx, metav1.ListOptions{})
    if err != nil {
        return nil, err
    }
    return services.Items, nil
}

func (m *KubernetesMonitor) DetectCrashLoopBackOff(ctx context.Context, namespace string) ([]model.Incident, error) {
    pods, err := m.ListPods(ctx, namespace)
    if err != nil {
        return nil, err
    }

    incidents := make([]model.Incident, 0)
    for _, pod := range pods {
        if incident := m.buildCrashLoopIncident(pod); incident != nil {
            incidents = append(incidents, *incident)
        }
    }
    return incidents, nil
}

func (m *KubernetesMonitor) DetectUnhealthyPods(ctx context.Context, namespace string) ([]model.Incident, error) {
    pods, err := m.ListPods(ctx, namespace)
    if err != nil {
        return nil, err
    }

    incidents := make([]model.Incident, 0)
    for _, pod := range pods {
        if isCrashLoopBackOff(pod) {
            continue
        }
        if incident := m.buildUnhealthyPodIncident(pod); incident != nil {
            incidents = append(incidents, *incident)
        }
    }
    return incidents, nil
}

func (m *KubernetesMonitor) DetectDeploymentReplicaMismatches(ctx context.Context, namespace string) ([]model.Incident, error) {
    deployments, err := m.ListDeployments(ctx, namespace)
    if err != nil {
        return nil, err
    }

    incidents := make([]model.Incident, 0)
    for _, deploy := range deployments {
        if incident := m.buildDeploymentReplicaIncident(deploy); incident != nil {
            incidents = append(incidents, *incident)
        }
    }
    return incidents, nil
}

func (m *KubernetesMonitor) buildCrashLoopIncident(pod corev1.Pod) *model.Incident {
    for _, status := range pod.Status.ContainerStatuses {
        if status.State.Waiting != nil && strings.Contains(status.State.Waiting.Reason, "CrashLoopBackOff") {
            return &model.Incident{
                ID:                  fmt.Sprintf("pod-%s-%s", pod.Namespace, pod.Name),
                Namespace:           pod.Namespace,
                Workload:            pod.Name,
                Type:                "CrashLoopBackOff",
                Summary:             "Pod is stuck in CrashLoopBackOff",
                Details:             fmt.Sprintf("container=%s reason=%s message=%s", status.Name, status.State.Waiting.Reason, status.State.Waiting.Message),
                Severity:            model.SeverityHigh,
                SeverityScore:       85,
                SuggestedRemediation: "Inspect pod logs and restart the container after fixing the crash loop.",
                FirstSeen:           time.Now(),
                LastSeen:            time.Now(),
                Source:              "kubernetes-monitor",
            }
        }
    }
    return nil
}

func (m *KubernetesMonitor) buildUnhealthyPodIncident(pod corev1.Pod) *model.Incident {
    if pod.Status.Phase == corev1.PodFailed || pod.Status.Phase == corev1.PodUnknown {
        return &model.Incident{
            ID:                  fmt.Sprintf("pod-%s-%s", pod.Namespace, pod.Name),
            Namespace:           pod.Namespace,
            Workload:            pod.Name,
            Type:                "UnhealthyPod",
            Summary:             "Pod is not healthy",
            Details:             fmt.Sprintf("phase=%s ready=%t restartCount=%d", pod.Status.Phase, isPodReady(pod), totalRestartCount(pod)),
            Severity:            model.SeverityHigh,
            SeverityScore:       70,
            SuggestedRemediation: "Evaluate container readiness and pod conditions to restore pod health.",
            FirstSeen:           time.Now(),
            LastSeen:            time.Now(),
            Source:              "kubernetes-monitor",
        }
    }

    if !isPodReady(pod) {
        return &model.Incident{
            ID:                  fmt.Sprintf("pod-%s-%s", pod.Namespace, pod.Name),
            Namespace:           pod.Namespace,
            Workload:            pod.Name,
            Type:                "UnhealthyPod",
            Summary:             "Pod readiness check failed",
            Details:             fmt.Sprintf("ready=%t restartCount=%d", isPodReady(pod), totalRestartCount(pod)),
            Severity:            model.SeverityMedium,
            SeverityScore:       55,
            SuggestedRemediation: "Review pod conditions and health checks to restore readiness.",
            FirstSeen:           time.Now(),
            LastSeen:            time.Now(),
            Source:              "kubernetes-monitor",
        }
    }
    return nil
}

func (m *KubernetesMonitor) buildDeploymentReplicaIncident(deploy appsv1.Deployment) *model.Incident {
    if deploy.Status.ReadyReplicas < deploy.Status.Replicas {
        return &model.Incident{
            ID:                  fmt.Sprintf("deployment-%s-%s", deploy.Namespace, deploy.Name),
            Namespace:           deploy.Namespace,
            Workload:            deploy.Name,
            Type:                "DeploymentReplicaMismatch",
            Summary:             "Deployment replica mismatch detected",
            Details:             fmt.Sprintf("desired=%d ready=%d available=%d", deploy.Status.Replicas, deploy.Status.ReadyReplicas, deploy.Status.AvailableReplicas),
            Severity:            model.SeverityMedium,
            SeverityScore:       60,
            SuggestedRemediation: "Inspect deployment status and adjust replica configuration or pod health checks.",
            FirstSeen:           time.Now(),
            LastSeen:            time.Now(),
            Source:              "kubernetes-monitor",
        }
    }
    return nil
}

func isPodReady(pod corev1.Pod) bool {
    for _, condition := range pod.Status.Conditions {
        if condition.Type == corev1.PodReady && condition.Status == corev1.ConditionTrue {
            return true
        }
    }
    return false
}

func totalRestartCount(pod corev1.Pod) int32 {
    var restarts int32
    for _, status := range pod.Status.ContainerStatuses {
        restarts += status.RestartCount
    }
    return restarts
}

func isCrashLoopBackOff(pod corev1.Pod) bool {
    for _, status := range pod.Status.ContainerStatuses {
        if status.State.Waiting != nil && strings.Contains(status.State.Waiting.Reason, "CrashLoopBackOff") {
            return true
        }
    }
    return false
}
