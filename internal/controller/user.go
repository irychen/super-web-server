package controller

import (
	"super-web-server/internal/ctx"
	"super-web-server/internal/dto"
	"super-web-server/internal/exception"
	"super-web-server/internal/service"
	"super-web-server/pkg/logger"

	"github.com/gin-gonic/gin"
)

type UserController interface {
	LoginByEmail(gtx *gin.Context)
	Info(gtx *gin.Context)
}

type userController struct {
	userService service.UserService
	logger      *logger.Logger
}

func NewUserController(userService service.UserService, logger *logger.Logger) UserController {
	logger.Info("NewUserController initialized successfully")
	return &userController{
		userService: userService,
		logger:      logger,
	}
}

func (c *userController) LoginByEmail(gtx *gin.Context) {
	appCtx := ctx.NewAppCtx(gtx)
	var req dto.UserLoginByEmailReqDTO
	if err := appCtx.ShouldBind(&req); err != nil {
		appCtx.ToError(exception.ExceptionInvalidParam.AppendDetails(*err...))
		return
	}
	data, err := c.userService.LoginByEmail(gtx, req)
	if err != nil {
		appCtx.ToError(err)
		return
	}
	appCtx.ToSuccess(data)
}

func (c *userController) Info(gtx *gin.Context) {
	appCtx := ctx.NewAppCtx(gtx)
	userUniqueID, err := appCtx.GetUserUniqueID()
	if err != nil {
		appCtx.ToError(exception.ExceptionUnauthorized.AppendDetails(err.Error()))
		return
	}
	user, ex := c.userService.GetUserByUniqueID(gtx, userUniqueID)
	if ex != nil {
		appCtx.ToError(ex)
		return
	}
	appCtx.ToSuccess(user)
}
