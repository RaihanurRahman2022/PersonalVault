package middleware

import (
	"time"

	"github.com/RaihanurRahman2022/PersonalVault/internal/monitoring"
	"github.com/gin-gonic/gin"
)

func MonitoringMiddleware(metrics *monitoring.Metrics) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()

		// Increment active connections
		metrics.IncrementActiveConnections()
		defer metrics.DecrementActiveConnections()

		// Process request
		ctx.Next()

		// Record metrics
		duration := time.Since(start)
		metrics.RecordRequest(
			ctx.Request.Method,
			ctx.FullPath(),
			ctx.Writer.Status(),
			duration,
		)
	}
}
func OpenTelemetryMiddleware() gin.HandlerFunc {
	return otelgin.Middleware("personal-vault-server")
}
