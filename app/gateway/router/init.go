package router

import (
	"github.com/gin-gonic/gin"
	"star/app/gateway/httpHandler"
)

func Setup() *gin.Engine {
	v := gin.New()
	// v.Use(logger.GinLogger(), logger.GinRecovery(true))

	// 用户相关路由
	v.GET("/signup/send", httpHandler.SendSetupHandler)
	v.POST("/signup", httpHandler.SignupHandler)
	v1 := v.Group("/login")
	{
		v1.POST("/password", httpHandler.LoginHandler)
		v1.GET("/captcha/send", httpHandler.SendLoginHandler)
		v1.POST("/captcha", httpHandler.LoginWithCaptcha)
	}

	// 评论相关路由
	v.POST("/comment", httpHandler.PostComment)
	v.DELETE("/comment/:id", httpHandler.DeleteComment)
	v.GET("/comments", httpHandler.GetComments)
	v.POST("/comment/star/:id", httpHandler.StarComment)

	return v
}
