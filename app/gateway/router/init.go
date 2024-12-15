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
	v1 := v.Group("/account")
	{
		v1.POST("/register", httpHandler.SignupHandler)
		v1.POST("/checkCode", httpHandler.GetCaptchaHandler)
		v1.POST("/login", httpHandler.LoginHandler)
		v1.POST("/autoLogin", httpHandler.AutoLoginHandler)
		v1.POST("/send", httpHandler.SendSetupHandler)
	}
	v.POST("/refreshToken", httpHandler.RefreshTokenHandler)
	v2 := v.Group("/admin")
	{
		v2.POST("/account/checkCode", httpHandler.GetCaptchaHandler)
		v2.POST("/account/login", httpHandler.LoginAdminHandler)

		v3 := v2.Use(middleware.AdminAuthHandler)
		{
			v3.POST("/category/loadCategory", httpHandler.LoadCategoryListHandler)
			v3.POST("/category/delCategory", httpHandler.DelCategoryHandler)
			v3.POST("/category/saveCategory", httpHandler.SaveCategoryHandler)
			v3.POST("/category/changeSort", httpHandler.ChangeSortHandler)
			v3.POST("/file/uploadImage", httpHandler.FileUploadHandler)
		}
	}
	v.POST("/category/loadAllCategory", httpHandler.LoadCategoryListHandler)
	v.POST("file/preUploadVideo", middleware.JWTAuthHandler, httpHandler.PreUploadVideosHandler)

	//

	//使用登录鉴权中间件
	//v.Use(middleware.JWTAuthHandler)
	//{
	//	// 评论相关路由
	//	v.POST("/comment", httpHandler.PostComment)
	//	v.DELETE("/comment/:id", httpHandler.DeleteComment)
	//	v.GET("/comments", httpHandler.GetComments)
	//
	//	v.GET("/whisper", httpHandler.ListMessageCountHandler)
	//	v.GET("/whisper/:userId", httpHandler.SendPrivateMessageHandler)
	//}

	return v
}
