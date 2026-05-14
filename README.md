# KubeWatch AI

AI-powered Kubernetes incident monitoring and root-cause analysis backend.

## What it includes

- Backend scaffold in Go with Gin.
- Modular folder organization under `backend/internal`.
- Kubernetes service structure using `client-go`.
- Prometheus metrics integration and health endpoint.
- Incident model and placeholder REST APIs.

## Run locally

```bash
cd backend
go run ./cmd/server
```

## Notes

The frontend is not generated yet. This repository currently contains the Go backend scaffold only.
