package router

import (
	"github.com/gin-gonic/gin"
	"star/app/gateway/httpHandler"
	"star/app/gateway/middleware"
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
	//使用登录鉴权中间件
	v.Use(middleware.JWTAuthHandler)
	{
		// 评论相关路由
		v.POST("/comment", httpHandler.PostComment)
		v.DELETE("/comment/:id", httpHandler.DeleteComment)
		v.GET("/comments", httpHandler.GetComments)
		v.POST("/comment/star/:id", httpHandler.StarComment)

		//v2 := v.Group("/community")
		//{
		//	//v2.POST()
		//
		//}

		v.POST("/community/:communityId/chat", httpHandler.ChatHandler)
		v.GET("/websocket", httpHandler.ListMessageCountHandler)
	}

	return v
}
