package analysis

import (
    "fmt"
    "strings"
    "time"

    "kubewatch-ai/internal/core/model"
)

type IncidentCategory string

const (
    CategoryCrashLoopBackOff   IncidentCategory = "crash-loop-backoff"
    CategoryImagePullBackOff   IncidentCategory = "image-pull-backoff"
    CategorySchedulingFailure  IncidentCategory = "scheduling-failure"
    CategoryUnhealthyDeployment IncidentCategory = "unhealthy-deployment"
    CategoryHighRestartCount   IncidentCategory = "high-restart-count"
    CategoryNodePressure       IncidentCategory = "node-pressure"
    CategoryGeneric            IncidentCategory = "generic"
)

type IncidentAnalyzer struct{}

func NewIncidentAnalyzer() *IncidentAnalyzer {
    return &IncidentAnalyzer{}
}

func (a *IncidentAnalyzer) AnalyzeCrashLoopBackOff(podName, namespace, container, reason, message string, restartCount int32) model.Incident {
    details := fmt.Sprintf("container=%s reason=%s message=%s restarts=%d", container, reason, message, restartCount)
    return a.buildIncident(podName, namespace, "CrashLoopBackOff", CategoryCrashLoopBackOff, "CrashLoopBackOff detected", details, model.SeverityHigh, 88, a.suggestionCrashLoopBackOff(reason, message), "Pod container repeatedly fails startup.")
}

func (a *IncidentAnalyzer) AnalyzeImagePullBackOff(podName, namespace, container, image, reason, message string) model.Incident {
    details := fmt.Sprintf("container=%s image=%s reason=%s message=%s", container, image, reason, message)
    return a.buildIncident(podName, namespace, "ImagePullBackOff", CategoryImagePullBackOff, "ImagePullBackOff detected", details, model.SeverityHigh, 82, a.suggestionImagePullBackOff(image, reason), "Pod cannot pull container image.")
}

func (a *IncidentAnalyzer) AnalyzeSchedulingFailure(podName, namespace, nodeSelector, message string) model.Incident {
    details := fmt.Sprintf("nodeSelector=%s message=%s", nodeSelector, message)
    return a.buildIncident(podName, namespace, "FailedScheduling", CategorySchedulingFailure, "Pod scheduling failed", details, model.SeverityHigh, 78, a.suggestionSchedulingFailure(nodeSelector, message), "Pod could not be placed on a node.")
}

func (a *IncidentAnalyzer) AnalyzeUnhealthyPod(podName, namespace, details string) model.Incident {
    return a.buildIncident(podName, namespace, "UnhealthyPod", CategoryUnhealthyDeployment, "Unhealthy pod state", details, model.SeverityHigh, 70, "Inspect pod and container health conditions, then resolve readiness or failure conditions.", "Pod is failing readiness or liveness checks.")
}

func (a *IncidentAnalyzer) AnalyzeUnhealthyDeployment(deploymentName, namespace, details string) model.Incident {
    return a.buildIncident(deploymentName, namespace, "UnhealthyDeployment", CategoryUnhealthyDeployment, "Unhealthy deployment state", details, model.SeverityMedium, 60, a.suggestionUnhealthyDeployment(details), "Deployment is not meeting readiness or availability expectations.")
}

func (a *IncidentAnalyzer) AnalyzeHighRestartCount(podName, namespace string, restartCount int32) model.Incident {
    details := fmt.Sprintf("restartCount=%d", restartCount)
    return a.buildIncident(podName, namespace, "HighRestartCount", CategoryHighRestartCount, "Pod exhibits high container restart count", details, model.SeverityMedium, 62, a.suggestionHighRestartCount(restartCount), "A container has restarted frequently.")
}

func (a *IncidentAnalyzer) AnalyzeNodePressure(nodeName, pressureType, message string) model.Incident {
    details := fmt.Sprintf("pressureType=%s message=%s", pressureType, message)
    return a.buildIncident(nodeName, "cluster", "NodePressure", CategoryNodePressure, "Node pressure condition detected", details, model.SeverityHigh, 75, a.suggestionNodePressure(pressureType), "Node resource pressure may affect workloads.")
}

func (a *IncidentAnalyzer) buildIncident(workload, namespace, incidentType string, category IncidentCategory, summary, details string, severity model.Severity, score int, remediation, rootCause string) model.Incident {
    return model.Incident{
        ID:                   fmt.Sprintf("%s-%s-%s", strings.ToLower(string(category)), namespace, workload),
        Namespace:            namespace,
        Workload:             workload,
        Type:                 incidentType,
        Category:             model.IncidentCategory(category),
        Summary:              summary,
        Details:              details,
        Severity:             severity,
        SeverityScore:        score,
        SuggestedRemediation: remediation,
        RootCause:            rootCause,
        FirstSeen:            time.Now(),
        LastSeen:             time.Now(),
        Source:               "analysis-engine",
    }
}

func (a *IncidentAnalyzer) suggestionCrashLoopBackOff(reason, message string) string {
    if strings.Contains(reason, "CrashLoopBackOff") {
        return "Inspect the failing container command and application logs. Fix the startup issue or resource constraints, then redeploy."
    }
    return "Check the pod logs for startup failures and evaluate readiness probes."
}

func (a *IncidentAnalyzer) suggestionImagePullBackOff(image, reason string) string {
    if strings.Contains(reason, "ErrImagePull") || strings.Contains(reason, "ImagePullBackOff") {
        return fmt.Sprintf("Verify image name '%s' and registry credentials. Ensure the image exists and the node can reach the registry.", image)
    }
    return "Check image repository access and authentication settings."
}

func (a *IncidentAnalyzer) suggestionSchedulingFailure(nodeSelector, message string) string {
    if nodeSelector != "" {
        return fmt.Sprintf("Review node selector and taint/toleration rules, and confirm available nodes match '%s'.", nodeSelector)
    }
    return fmt.Sprintf("Inspect cluster resources and scheduler events: %s", message)
}

func (a *IncidentAnalyzer) suggestionUnhealthyDeployment(details string) string {
    return fmt.Sprintf("Review deployment status, pod health checks, and rollout history. Address the underlying issue: %s", details)
}

func (a *IncidentAnalyzer) suggestionHighRestartCount(count int32) string {
    if count > 10 {
        return "A high restart count indicates instability. Check crash loops, OOM kills, and liveness probes."
    }
    return "Investigate container health and startup behavior to reduce restarts."
}

func (a *IncidentAnalyzer) suggestionNodePressure(pressureType string) string {
    return fmt.Sprintf("Node pressure detected on %s. Evaluate node utilization, free resources, and consider scaling the cluster.", pressureType)
}
