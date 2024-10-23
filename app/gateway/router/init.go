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
		v1.POST("/captcha", httpHandler.LoginWithCaptchaHandler)
	}

	//

	v2 := v.Group("/community")
	{
		v2.POST("/create", middleware.JWTAuthHandler, httpHandler.CreateCommunityHandler)
		v2.POST("/:id/follow", middleware.JWTAuthHandler, httpHandler.FollowCommunityHandler)
		v2.POST("/:id/unFollow", middleware.JWTAuthHandler, httpHandler.UnFollowCommunityHandler)
		v2.POST("publish", middleware.JWTAuthHandler, httpHandler.CreatePostHandler)

	}
	v.POST("/like", middleware.JWTAuthHandler, httpHandler.LikeActionHandler)
	//使用登录鉴权中间件
	v.Use(middleware.JWTAuthHandler)
	{
		// 评论相关路由
		v.POST("/comment", httpHandler.PostComment)
		v.DELETE("/comment/:id", httpHandler.DeleteComment)
		v.GET("/comments", httpHandler.GetComments)

		v.GET("/whisper", httpHandler.ListMessageCountHandler)
		v.GET("/whisper/:userId", httpHandler.SendPrivateMessageHandler)
	}

	return v
}
