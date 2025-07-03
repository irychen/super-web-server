package controller

import (
	"super-web-server/internal/ctx"
	"super-web-server/pkg/logger"
	"time"

	"github.com/gin-gonic/gin"
)

type HelloController interface {
	Hello(gtx *gin.Context)
}

type helloController struct {
	logger *logger.Logger
}

func NewHelloController(logger *logger.Logger) HelloController {
	logger.Info("NewHelloController initialized successfully")
	return &helloController{logger: logger}
}

func (c *helloController) Hello(gtx *gin.Context) {
	appCtx := ctx.NewAppCtx(gtx)
	time.Sleep(1 * time.Second)
	appCtx.ToSuccess("Hello, World!")
}
