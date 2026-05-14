package rest

import (
    "net/http"

    gorillaws "github.com/gorilla/websocket"
    "github.com/gin-gonic/gin"
    ws "kubewatch-ai/internal/adapter/websocket"
    "kubewatch-ai/internal/core/service"
)

func OverviewHandler(c *gin.Context, service *service.IncidentService) {
    overview, err := service.GetClusterOverview(c.Request.Context())
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, overview)
}

func NamespacesHandler(c *gin.Context, service *service.IncidentService) {
    namespaces, err := service.GetNamespaces(c.Request.Context())
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"namespaces": namespaces})
}

func IncidentListHandler(c *gin.Context, service *service.IncidentService) {
    incidents, err := service.GetIncidents(c.Request.Context())
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"incidents": incidents})
}

func WebSocketHandler(c *gin.Context, hub *ws.Hub) {
    upgrader := gorillaws.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
    conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
    if err != nil {
        c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    client := ws.NewClient(conn)
    hub.Register <- client
}
