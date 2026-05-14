package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	gorillaws "github.com/gorilla/websocket"
	"k8s.io/api/core/v1"
	ws "kubewatch-ai/internal/adapter/websocket"
	"kubewatch-ai/internal/core/service"
)

type apiError struct {
	Error string `json:"error"`
}

type PodSummary struct {
	Name       string `json:"name"`
	Namespace  string `json:"namespace"`
	Phase      string `json:"phase"`
	NodeName   string `json:"nodeName"`
	Restarts   int32  `json:"restarts"`
	Containers int    `json:"containers"`
}

func respondSuccess(c *gin.Context, status int, payload interface{}) {
	c.JSON(status, payload)
}

func respondError(c *gin.Context, status int, err error) {
	c.JSON(status, apiError{Error: err.Error()})
}

func totalPodRestarts(pod v1.Pod) int32 {
	var restarts int32
	for _, status := range pod.Status.ContainerStatuses {
		restarts += status.RestartCount
	}
	return restarts
}

func OverviewHandler(c *gin.Context, service *service.IncidentService) {
	overview, err := service.GetClusterOverview(c.Request.Context())
	if err != nil {
		respondError(c, http.StatusInternalServerError, err)
		return
	}
	respondSuccess(c, http.StatusOK, overview)
}

func HealthHandler(c *gin.Context) {
	respondSuccess(c, http.StatusOK, gin.H{"status": "ok"})
}

func NamespacesHandler(c *gin.Context, service *service.IncidentService) {
	namespaces, err := service.GetNamespaces(c.Request.Context())
	if err != nil {
		respondError(c, http.StatusInternalServerError, err)
		return
	}
	respondSuccess(c, http.StatusOK, gin.H{"namespaces": namespaces})
}

func PodsHandler(c *gin.Context, service *service.IncidentService) {
	pods, err := service.GetPods(c.Request.Context())
	if err != nil {
		respondError(c, http.StatusInternalServerError, err)
		return
	}

	summaries := make([]PodSummary, 0, len(pods))
	for _, pod := range pods {
		summaries = append(summaries, PodSummary{
			Name:       pod.Name,
			Namespace:  pod.Namespace,
			Phase:      string(pod.Status.Phase),
			NodeName:   pod.Spec.NodeName,
			Restarts:   totalPodRestarts(pod),
			Containers: len(pod.Spec.Containers),
		})
	}
	respondSuccess(c, http.StatusOK, gin.H{"pods": summaries})
}

func IncidentListHandler(c *gin.Context, service *service.IncidentService) {
	incidents, err := service.GetIncidents(c.Request.Context())
	if err != nil {
		respondError(c, http.StatusInternalServerError, err)
		return
	}
	respondSuccess(c, http.StatusOK, gin.H{"incidents": incidents})
}

func ClusterHealthHandler(c *gin.Context, service *service.IncidentService) {
	health, err := service.GetClusterHealth(c.Request.Context())
	if err != nil {
		respondError(c, http.StatusInternalServerError, err)
		return
	}
	respondSuccess(c, http.StatusOK, health)
}

func WebSocketHandler(c *gin.Context, hub *ws.Hub) {
	upgrader := gorillaws.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, apiError{Error: err.Error()})
		return
	}
	client := ws.NewClient(conn)
	hub.Register <- client
}
