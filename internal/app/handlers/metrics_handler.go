package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type MetricsHandler struct{}

func NewMetricsHandler() *MetricsHandler {
	return &MetricsHandler{}
}

func (h *MetricsHandler) GetMetrics(c *gin.Context) {
	promhttp.Handler().ServeHTTP(c.Writer, c.Request)
}

func (h *MetricsHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"service":   "personal-vault-server",
		"timestamp": "2024-01-01T00:00:00Z",
	})
}
