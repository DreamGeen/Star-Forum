package router

import (
	"github.com/gin-gonic/gin"

	"star/app/gateway/httpHandler"
)

func Setup() *gin.Engine {
	v := gin.New()
	//v.Use(logger.GinLogger(), logger.GinRecovery(true))
	//注册
	{
		v.GET("/signup/send", httpHandler.SendSetupHandler)
		v.POST("/signup", httpHandler.SignupHandler)
	}
	//登录
	v1 := v.Group("/login")
	{
		v1.POST("/password", httpHandler.LoginHandler)
		v1.GET("/captcha/send", httpHandler.SendLoginHandler)
		v1.POST("/captcha", httpHandler.LoginWithCaptcha)
	}

	return v
}
