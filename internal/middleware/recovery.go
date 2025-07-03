package middleware

import (
	"super-web-server/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Recovery() gin.HandlerFunc {
	ginLogger := logger.GetModuleLogger("gin")
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				ginLogger.Error("gin panic recovered",
					zap.Any("error", err),
					zap.String("request", c.Request.Method+" "+c.Request.URL.Path),
				)
				c.AbortWithStatus(500)
			}
		}()
		c.Next()
	}
}
