package app

import (
	"fmt"
	"net/http"
	"super-web-server/internal/middleware"
	"super-web-server/internal/types"

	"github.com/gin-gonic/gin"
)

func (a *App) InitEngineAndServer() {
	var serverConfig = a.config.Server
	if a.config.Mode != types.ServerModeDev {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	a.engine = gin.New()
	a.engine.Use(middleware.Recovery())
	a.engine.Use(middleware.Logger())
	a.engine.Use(middleware.CORS())
	a.engine.Use(middleware.Translations())

	server := &http.Server{
		Addr:           fmt.Sprintf(":%d", serverConfig.Port),
		Handler:        a.engine,
		ReadTimeout:    serverConfig.ReadTimeout,
		WriteTimeout:   serverConfig.WriteTimeout,
		MaxHeaderBytes: serverConfig.MaxHeaderBytes,
	}
	a.server = server
}
