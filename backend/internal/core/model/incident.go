package model

import "time"

type Severity string

const (
    SeverityCritical Severity = "critical"
    SeverityHigh     Severity = "high"
    SeverityMedium   Severity = "medium"
    SeverityLow      Severity = "low"
)

type Incident struct {
    ID                  string    `json:"id"`
    Namespace           string    `json:"namespace"`
    Workload            string    `json:"workload"`
    Type                string    `json:"type"`
    Summary             string    `json:"summary"`
    Details             string    `json:"details"`
    Severity            Severity  `json:"severity"`
    SeverityScore       int       `json:"severityScore"`
    SuggestedRemediation string   `json:"suggestedRemediation"`
    FirstSeen           time.Time `json:"firstSeen"`
    LastSeen            time.Time `json:"lastSeen"`
    Source              string    `json:"source"`
}

type ClusterOverview struct {
    TotalNodes          int `json:"totalNodes"`
    ReadyNodes          int `json:"readyNodes"`
    TotalNamespaces     int `json:"totalNamespaces"`
    CrashLoopBackOff    int `json:"crashLoopBackOff"`
    UnhealthyDeployments int `json:"unhealthyDeployments"`
    ReplicaMismatch     int `json:"replicaMismatch"`
    FailedServices      int `json:"failedServices"`
    ActiveIncidents     int `json:"activeIncidents"`
}
