package ctx

import (
	"errors"
	"net/http"
	"super-web-server/internal/dto"
	"super-web-server/internal/exception"
	"super-web-server/internal/validator"

	"github.com/gin-gonic/gin"
)

type AppCtx struct {
	*gin.Context
	validator *validator.Validator
}

const (
	SUCCESS_CODE          = 0
	SUCCESS_MESSAGE       = "success"
	USER_UNIQUE_ID_KEY    = "user_unique_id"
	USER_UNIQUE_ROLES_KEY = "user_unique_roles"
)

func NewAppCtx(gtx *gin.Context) *AppCtx {
	return &AppCtx{
		Context:   gtx,
		validator: validator.NewValidator(gtx),
	}
}

func (c *AppCtx) ShouldBind(obj any) *[]string {
	return c.validator.ShouldBind(obj)
}

func (c *AppCtx) GetUserUniqueID() (int64, error) {
	id := c.GetInt64(USER_UNIQUE_ID_KEY)
	if id == 0 {
		return 0, errors.New("user unique id not found")
	}
	return id, nil
}

func (c *AppCtx) SetUserUniqueID(id int64) {
	c.Set(USER_UNIQUE_ID_KEY, id)
}

func (c *AppCtx) ToError(err *exception.Exception) {
	c.JSON(err.StatusCode, gin.H{
		"code":    err.Code,
		"message": err.Message,
		"details": err.Details,
	})
	c.Abort()
}

func (c *AppCtx) ToSuccess(data any) {
	c.JSON(http.StatusOK, gin.H{
		"code":    SUCCESS_CODE,
		"message": SUCCESS_MESSAGE,
		"data":    data,
	})
	c.Abort()
}

func (c *AppCtx) ToSuccessPageList(data any, total int64, pagination *dto.Pagination) {
	c.JSON(http.StatusOK, gin.H{
		"code":    SUCCESS_CODE,
		"message": SUCCESS_MESSAGE,
		"data": gin.H{
			"list":     data,
			"total":    total,
			"pageSize": pagination.Value().PageSize,
			"page":     pagination.Value().Page,
		},
	})
	c.Abort()
}
