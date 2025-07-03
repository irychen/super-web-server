package middleware

import (
	"fmt"
	"slices"
	"super-web-server/internal/ctx"
	"super-web-server/internal/exception"
	"super-web-server/internal/model"
	"super-web-server/internal/service"

	"github.com/gin-gonic/gin"
)

type RoleCheck struct {
	service service.Service
}

func NewRoleCheck(service service.Service) *RoleCheck {
	return &RoleCheck{
		service: service,
	}
}

func (r *RoleCheck) RoleCheckAll(requiredRoles ...model.UserRoleEnum) gin.HandlerFunc {
	return func(c *gin.Context) {
		appCtx := ctx.NewAppCtx(c)
		userUniqueID, err := appCtx.GetUserUniqueID()
		if err != nil {
			appCtx.ToError(exception.ExceptionUnauthorized.AppendDetails(err.Error()))
			return
		}
		userRoles, ex := r.service.User().GetUserCachedRolesByUniqueID(c, userUniqueID)
		if ex != nil {
			appCtx.ToError(ex)
			return
		}

		var userRoleCodes []model.UserRoleEnum
		for _, role := range userRoles {
			userRoleCodes = append(userRoleCodes, role.Code)
		}

		for _, requiredRole := range requiredRoles {
			if !slices.Contains(userRoleCodes, requiredRole) {
				appCtx.ToError(exception.ExceptionUnauthorized.AppendDetails(fmt.Sprintf("user role %s not in %v", requiredRole, userRoleCodes)))
				return
			}
		}

		c.Next()
	}
}

func (r *RoleCheck) RoleCheckAny(requiredRoles ...model.UserRoleEnum) gin.HandlerFunc {
	return func(c *gin.Context) {
		appCtx := ctx.NewAppCtx(c)
		userUniqueID, err := appCtx.GetUserUniqueID()
		if err != nil {
			appCtx.ToError(exception.ExceptionUnauthorized.AppendDetails(err.Error()))
			return
		}
		userRoles, ex := r.service.User().GetUserCachedRolesByUniqueID(c, userUniqueID)
		if ex != nil {
			appCtx.ToError(ex)
			return
		}

		var userRoleCodes []model.UserRoleEnum
		for _, role := range userRoles {
			userRoleCodes = append(userRoleCodes, role.Code)
		}

		for _, requiredRole := range requiredRoles {
			if slices.Contains(userRoleCodes, requiredRole) {
				c.Next()
				return
			}
		}

		appCtx.ToError(exception.ExceptionUnauthorized.AppendDetails(fmt.Sprintf("user role %v not in any of %v", userRoleCodes, requiredRoles)))
	}
}
