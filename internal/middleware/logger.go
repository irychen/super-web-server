package middleware

import (
	"fmt"
	"super-web-server/pkg/logger"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Logger() gin.HandlerFunc {
	ginLogger := logger.GetModuleLogger("gin")
	return func(c *gin.Context) {
		startTime := time.Now()
		c.Next()

		latency := time.Since(startTime)
		latencyMs := float64(latency.Microseconds()) / 1000.0

		fields := []zap.Field{
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("query", c.Request.URL.RawQuery),
			zap.String("ip", c.ClientIP()),
			zap.String("user_agent", c.Request.UserAgent()),
			zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
			zap.String("latency", fmt.Sprintf("%.2fms", latencyMs)),
		}

		status := c.Writer.Status()
		switch {
		case status >= 500:
			ginLogger.Error("Internal Server Error", fields...)
		case status >= 400:
			ginLogger.Warn("Failed", fields...)
		case status >= 300:
			ginLogger.Warn("Redirect", fields...)
		default:
			ginLogger.Info("Success", fields...)
		}
	}
}
