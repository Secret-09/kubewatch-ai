package main

import (
    "context"
    "log"
    "os/signal"
    "syscall"
    "time"

    "github.com/gin-gonic/gin"
    "kubewatch-ai/internal/adapter/rest"
    "kubewatch-ai/internal/adapter/websocket"
    "kubewatch-ai/internal/core/service"
    "kubewatch-ai/internal/infrastructure/k8s"
    "kubewatch-ai/internal/infrastructure/monitoring"
)

func main() {
    gin.SetMode(gin.ReleaseMode)

    k8sClient, err := k8s.NewKubernetesClient()
    if err != nil {
        log.Fatalf("failed to build kubernetes client: %v", err)
    }

    metrics := monitoring.NewPrometheusMetrics()
    hub := websocket.NewHub()
    incidentService := service.NewIncidentService(k8sClient, metrics, hub)

    ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
    defer stop()

    go hub.Run(ctx)
    go incidentService.StartPolling(ctx)

    engine := rest.NewRouter(incidentService, metrics, hub)
    addr := ":8080"
    log.Printf("starting KubeWatch AI backend on %s", addr)

    serverErr := make(chan error, 1)
    go func() {
        serverErr <- engine.Run(addr)
    }()

    select {
    case <-ctx.Done():
        log.Println("shutdown requested")
        return
    case err := <-serverErr:
        log.Fatalf("server failed: %v", err)
    case <-time.After(30 * time.Second):
        log.Println("backend startup completed")
    }
}
