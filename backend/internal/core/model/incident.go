package model

import "time"

type Severity string

const (
	SeverityCritical Severity = "critical"
	SeverityHigh     Severity = "high"
	SeverityMedium   Severity = "medium"
	SeverityLow      Severity = "low"
)

type IncidentCategory string

const (
	IncidentCategoryCrashLoopBackOff    IncidentCategory = "crash-loop-backoff"
	IncidentCategoryImagePullBackOff    IncidentCategory = "image-pull-backoff"
	IncidentCategorySchedulingFailure   IncidentCategory = "scheduling-failure"
	IncidentCategoryUnhealthyDeployment IncidentCategory = "unhealthy-deployment"
	IncidentCategoryHighRestartCount    IncidentCategory = "high-restart-count"
	IncidentCategoryNodePressure        IncidentCategory = "node-pressure"
	IncidentCategoryGeneric             IncidentCategory = "generic"
)

type Incident struct {
	ID                   string           `json:"id"`
	Namespace            string           `json:"namespace"`
	Workload             string           `json:"workload"`
	Type                 string           `json:"type"`
	Category             IncidentCategory `json:"category"`
	Summary              string           `json:"summary"`
	Details              string           `json:"details"`
	Severity             Severity         `json:"severity"`
	SeverityScore        int              `json:"severityScore"`
	SuggestedRemediation string           `json:"suggestedRemediation"`
	RootCause            string           `json:"rootCause,omitempty"`
	FirstSeen            time.Time        `json:"firstSeen"`
	LastSeen             time.Time        `json:"lastSeen"`
	Source               string           `json:"source"`
}

type ClusterOverview struct {
	TotalNodes           int `json:"totalNodes"`
	ReadyNodes           int `json:"readyNodes"`
	TotalNamespaces      int `json:"totalNamespaces"`
	CrashLoopBackOff     int `json:"crashLoopBackOff"`
	UnhealthyDeployments int `json:"unhealthyDeployments"`
	ReplicaMismatch      int `json:"replicaMismatch"`
	FailedServices       int `json:"failedServices"`
	ActiveIncidents      int `json:"activeIncidents"`
}

type ClusterHealth struct {
	Status          string          `json:"status"`
	Overview        ClusterOverview `json:"overview"`
	TotalPods       int             `json:"totalPods"`
	TotalNamespaces int             `json:"totalNamespaces"`
	IncidentCount   int             `json:"incidentCount"`
	UnhealthyPods   int             `json:"unhealthyPods"`
	LastUpdated     time.Time       `json:"lastUpdated"`
}
