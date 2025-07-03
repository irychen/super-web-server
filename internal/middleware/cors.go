package middleware

import (
	"slices"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type CORSConfig struct {
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	ExposeHeaders    []string
	AllowCredentials bool
	MaxAge           int
}

func DefaultCORSConfig() CORSConfig {

	// 需要暴露给客户端的响应headers（主要用于ExposeHeaders）
	exposedHeaders := []string{
		"New-Token",
		"token",
		"x-token",
		"x-request-id",
		"accept",
		"origin",
		"Cache-Control",
		"X-Requested-With",
	}

	// 基础的HTTP协议headers
	basicHttpHeaders := []string{
		"Content-Type",
		"Content-Length",
		"Accept-Encoding",
		"X-CSRF-Token",
		"Authorization",
	}

	// 合并 headers
	allowHeaders := append(basicHttpHeaders, exposedHeaders...)

	return CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     allowHeaders,
		ExposeHeaders:    exposedHeaders,
		AllowCredentials: false,
		MaxAge:           172800,
	}
}

func CORSWithConfig(config CORSConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		var allowOrigin = false

		origin := c.Request.Header.Get("Origin")

		if len(config.AllowOrigins) > 0 {
			if slices.Contains(config.AllowOrigins, "*") {
				c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
				allowOrigin = false
			} else if slices.Contains(config.AllowOrigins, origin) {
				c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
				allowOrigin = true
			} else {
				// Origin 不在允许列表中，不设置CORS头，让浏览器处理
				allowOrigin = false
			}
		}

		if config.AllowCredentials && allowOrigin {
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		}

		if len(config.AllowHeaders) > 0 {
			c.Writer.Header().Set("Access-Control-Allow-Headers", strings.Join(config.AllowHeaders, ", "))
		}

		if len(config.AllowMethods) > 0 {
			c.Writer.Header().Set("Access-Control-Allow-Methods", strings.Join(config.AllowMethods, ", "))
		}

		if len(config.ExposeHeaders) > 0 {
			c.Writer.Header().Set("Access-Control-Expose-Headers", strings.Join(config.ExposeHeaders, ", "))
		}

		if config.MaxAge > 0 {
			c.Writer.Header().Set("Access-Control-Max-Age", strconv.Itoa(config.MaxAge))
		}

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

func CORS() gin.HandlerFunc {
	return CORSWithConfig(DefaultCORSConfig())
}

// 如果需要支持凭证的CORS配置
func CORSWithCredentials(allowOrigins []string) gin.HandlerFunc {
	config := DefaultCORSConfig()
	config.AllowOrigins = allowOrigins
	config.AllowCredentials = true
	return CORSWithConfig(config)
}
