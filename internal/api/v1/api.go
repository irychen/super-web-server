package v1

import (
	"super-web-server/internal/controller"
	"super-web-server/internal/middleware"
	"super-web-server/internal/model"
	"super-web-server/pkg/jwt"

	"github.com/gin-gonic/gin"
)

func InitApi(router *gin.RouterGroup, controller controller.Controller, jwt *jwt.JWT, rc *middleware.RoleCheck) {
	router.GET("/hello", controller.Hello().Hello)

	user := router.Group("/user")
	{
		user.POST("/login-by-email", controller.User().LoginByEmail)
	}

	user.Use(jwt.JWT(), rc.RoleCheckAny(
		model.UserRoleCodeSuperAdmin,
		model.UserRoleCodeAdmin,
		model.UserRoleCodeUser,
	))
	{
		user.GET("/info", controller.User().Info)
	}
}
