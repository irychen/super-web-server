package validator

import (
	"errors"
	"io"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

type Validator struct {
	ctx   *gin.Context
	trans ut.Translator
}

func Init() {
	jsonTagFunc := func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" || name == "" {
			name = strings.SplitN(fld.Tag.Get("form"), ",", 2)[0]
		}
		if name == "-" {
			return ""
		}
		return name
	}
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(jsonTagFunc)
	}
}

func NewValidator(ctx *gin.Context) *Validator {
	var trans ut.Translator
	if t, ok := ctx.Get("trans"); ok {
		trans = t.(ut.Translator)
	}
	var v = &Validator{
		ctx:   ctx,
		trans: trans,
	}
	return v
}

func (v *Validator) ShouldBind(s any) *[]string {
	if err := v.ctx.ShouldBind(s); err != nil {
		switch {
		case errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) || err.Error() == "EOF":
			return &[]string{"EOF body is empty"}
		default:
			if validationErrors, ok := err.(validator.ValidationErrors); ok {
				var verrors = make([]string, 0, len(validationErrors))
				for _, err := range validationErrors {
					var dtoPath = strings.SplitN(err.Namespace(), ".", 2)
					if len(dtoPath) > 1 {
						var translate = strings.Replace(err.Translate(v.trans), err.Field(), "", 1)
						verrors = append(verrors, dtoPath[1]+translate)
					} else {
						verrors = append(verrors, err.Translate(v.trans))
					}
				}
				return &verrors
			}
			return &[]string{err.Error()}
		}
	}
	return nil
}
