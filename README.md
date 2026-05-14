# KubeWatch AI

KubeWatch AI is a cloud-native observability platform designed for Kubernetes environments. It combines cluster health monitoring, incident detection, and AI-driven analysis into a modular backend service that scales with modern infrastructure.

## Project Overview

KubeWatch AI provides a centralized backend for monitoring Kubernetes clusters, collecting pod and namespace data, detecting incidents, and exposing structured metrics for dashboard or alerting systems. The architecture is built for cloud-native deployments with a focus on API-first integration, observability, and extensibility.

## Architecture

- **Gin HTTP server** handles REST API traffic and health checks.
- **Kubernetes client-go** interacts with cluster resources like Pods, Namespaces, Deployments, Nodes, and Services.
- **Background polling service** continuously gathers cluster state and incident information.
- **Incident analysis pipeline** transforms raw Kubernetes events into actionable incident objects.
- **Prometheus metrics** expose observability data for monitoring and alerting.
- **WebSocket hub** enables real-time incident notifications.

## System Architecture

![KubeWatch AI Architecture](docs/architecture.png)


## Tech Stack

- Go
- Gin web framework
- Kubernetes `client-go`
- Prometheus metrics
- Gorilla WebSocket
- Kubernetes API

## Backend Modules

- `backend/cmd/server`: application entrypoint and HTTP server bootstrap.
- `backend/internal/adapter/rest`: Gin router and REST handlers.
- `backend/internal/adapter/websocket`: real-time notification hub.
- `backend/internal/core/service`: core business logic for incident polling, cluster overview, and monitoring.
- `backend/internal/core/model`: domain models for incidents and cluster health.
- `backend/internal/infrastructure/k8s`: Kubernetes client initialization and cluster access.
- `backend/internal/infrastructure/monitoring`: Prometheus metrics instrumentation.

## API Endpoints

- `GET /health` - basic service health status
- `GET /api/v1/pods` - list pod summaries across the cluster
- `GET /api/v1/namespaces` - list cluster namespaces
- `GET /api/v1/incidents` - list detected incidents
- `GET /api/v1/cluster-health` - aggregated cluster health overview

## Kubernetes Monitoring Features

- Enumerates namespaces, pods, deployments, nodes, and services.
- Tracks pod health and readiness state.
- Computes cluster resource summaries and health indicators.
- Exposes Prometheus-compatible metrics for integration into cloud-native monitoring stacks.

## Incident Analysis Features

- Detects CrashLoopBackOff conditions.
- Identifies unhealthy pods and deployment replica mismatches.
- Converts Kubernetes events into structured incidents with severity, summary, and diagnostics.
- Supports real-time alert delivery through WebSockets.

## Planned Frontend Dashboard

A future frontend dashboard will provide:

- Cluster health overview and incident timeline
- Namespace and pod exploration views
- Incident severity filtering and remediation guidance
- Live alert stream from WebSocket notifications
- Cloud-native UX for Kubernetes operators and SRE teams

## Future Roadmap

- Add a full-featured frontend dashboard with charts and incident timelines.
- Expand incident analysis to include node, service, and network issues.
- Add RBAC-aware cluster access and multi-cluster support.
- Implement alerting integrations for Slack, Teams, and PagerDuty.
- Add policy-based anomaly detection and predictive failure analysis.
- Support managed Kubernetes platforms and cloud-provider metadata.

## Setup Instructions

### Prerequisites

- Go 1.21+ installed
- Access to a Kubernetes cluster or local `kubectl` configuration
- `GOPATH` / module-aware Go environment

### Run locally

```bash
cd /workspaces/kubewatch-ai/backend
go run ./cmd/server
```

### Run in cluster

1. Build the backend container image.
2. Deploy to Kubernetes with appropriate RBAC permissions.
3. Configure service and ingress to expose the REST API.

## Folder Structure

```
/README.md
/backend
  /cmd/server
    main.go
  /internal
    /adapter
      /rest
        handlers.go
        router.go
      /websocket
        hub.go
    /core
      /analysis
        engine.go
      /model
        incident.go
      /service
        incident_service.go
        kubernetes_monitor.go
    /infrastructure
      /k8s
        client.go
      /monitoring
        metrics.go
/frontend
  /app
    globals.css
    layout.tsx
    page.tsx
  package.json
  tsconfig.json
  next.config.ts
```

## Notes

This repository is currently backend-first with a frontend scaffold present for future dashboard development. The design prioritizes cloud-native observability, API-driven operations, and Kubernetes-native incident detection.
