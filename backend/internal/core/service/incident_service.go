package service

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	corev1 "k8s.io/api/core/v1"

	"kubewatch-ai/internal/adapter/websocket"
	"kubewatch-ai/internal/core/model"
	"kubewatch-ai/internal/infrastructure/k8s"
	"kubewatch-ai/internal/infrastructure/monitoring"
)

type IncidentService struct {
	monitor         *KubernetesMonitor
	metrics         *monitoring.PrometheusMetrics
	hub             *websocket.Hub
	mu              sync.RWMutex
	incidents       []model.Incident
	clusterOverview model.ClusterOverview
}

func NewIncidentService(client *k8s.Client, metrics *monitoring.PrometheusMetrics, hub *websocket.Hub) *IncidentService {
	return &IncidentService{monitor: NewKubernetesMonitor(client), metrics: metrics, hub: hub}
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
	return s.monitor.ListNamespaces(ctx)
}

func (s *IncidentService) GetPods(ctx context.Context) ([]corev1.Pod, error) {
	return s.monitor.ListPods(ctx, "")
}

func (s *IncidentService) GetClusterHealth(ctx context.Context) (model.ClusterHealth, error) {
	s.mu.RLock()
	overview := s.clusterOverview
	s.mu.RUnlock()

	pods, err := s.monitor.ListPods(ctx, "")
	if err != nil {
		return model.ClusterHealth{}, err
	}

	namespaces, err := s.monitor.ListNamespaces(ctx)
	if err != nil {
		return model.ClusterHealth{}, err
	}

	status := "healthy"
	if overview.ActiveIncidents > 0 {
		status = "degraded"
	}

	return model.ClusterHealth{
		Status:          status,
		Overview:        overview,
		TotalPods:       len(pods),
		TotalNamespaces: len(namespaces),
		IncidentCount:   overview.ActiveIncidents,
		UnhealthyPods:   overview.CrashLoopBackOff + overview.UnhealthyDeployments,
		LastUpdated:     time.Now(),
	}, nil
}

func (s *IncidentService) GetIncidents(ctx context.Context) ([]model.Incident, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]model.Incident, len(s.incidents))
	copy(out, s.incidents)
	return out, nil
}

func (s *IncidentService) refreshCluster(ctx context.Context) {
	pods, _ := s.monitor.ListPods(ctx, "")
	deployments, _ := s.monitor.ListDeployments(ctx, "")
	nodes, _ := s.monitor.ListNodes(ctx)
	services, _ := s.monitor.ListServices(ctx, "")
	namespaces, _ := s.monitor.ListNamespaces(ctx)

	incidents := make([]model.Incident, 0)
	incidents = append(incidents, mustDetect(s.monitor.DetectCrashLoopBackOff(ctx, ""))...)
	incidents = append(incidents, mustDetect(s.monitor.DetectUnhealthyPods(ctx, ""))...)
	incidents = append(incidents, mustDetect(s.monitor.DetectDeploymentReplicaMismatches(ctx, ""))...)

	incidentCounters := map[string]int{}
	incidentNamespaces := map[string]struct{}{}
	for _, incident := range incidents {
		incidentCounters[incident.Type]++
		incidentNamespaces[incident.Namespace] = struct{}{}
	}

	overview := model.ClusterOverview{
		TotalNodes:           len(nodes),
		ReadyNodes:           countReadyNodes(nodes),
		TotalNamespaces:      len(namespaces),
		CrashLoopBackOff:     incidentCounters["CrashLoopBackOff"],
		UnhealthyDeployments: incidentCounters["UnhealthyPod"],
		ReplicaMismatch:      incidentCounters["DeploymentReplicaMismatch"],
		FailedServices:       incidentCounters["FailedService"],
		ActiveIncidents:      len(incidents),
	}

	s.mu.Lock()
	s.incidents = incidents
	s.clusterOverview = overview
	s.mu.Unlock()

	s.metrics.UpdateIncidentCounts(len(incidents))
	s.metrics.UpdateOverviewMetrics(len(nodes), len(pods), len(deployments), len(services))

	payload, _ := json.Marshal(wsPayload{Incidents: incidents, Timestamp: time.Now()})
	s.hub.Broadcast(payload)
}

type wsPayload struct {
	Incidents []model.Incident `json:"incidents"`
	Timestamp time.Time        `json:"timestamp"`
}

func mustDetect(incidents []model.Incident, err error) []model.Incident {
	if err != nil {
		return nil
	}
	return incidents
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
