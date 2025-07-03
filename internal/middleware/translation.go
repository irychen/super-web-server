package middleware

import (
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_trans "github.com/go-playground/validator/v10/translations/en"
	zh_trans "github.com/go-playground/validator/v10/translations/zh"
)

var (
	uni  *ut.UniversalTranslator
	once sync.Once
)

func initTranslator() {
	once.Do(func() {
		uni = ut.New(en.New(), zh.New())
		// 预先注册所有支持的语言翻译
		if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
			registerZhTrans(v)
			registerEnTrans(v)
		}
	})
}

func registerZhTrans(v *validator.Validate) {
	if trans, _ := uni.GetTranslator("zh"); trans != nil {
		_ = zh_trans.RegisterDefaultTranslations(v, trans)
	}
}

func registerEnTrans(v *validator.Validate) {
	if trans, _ := uni.GetTranslator("en"); trans != nil {
		_ = en_trans.RegisterDefaultTranslations(v, trans)
	}
}

func Translations() gin.HandlerFunc {
	initTranslator()
	return func(c *gin.Context) {
		locale := c.GetHeader("locale")
		if locale == "" {
			locale = "zh"
		}
		trans, _ := uni.GetTranslator(locale)
		c.Set("trans", trans)
		c.Next()
	}
}
