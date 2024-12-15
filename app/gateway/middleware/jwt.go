package middleware

import (
	"log"
	"star/app/constant/settings"
	"star/app/constant/str"
	"star/app/utils/jwt"
	"strings"

	"github.com/gin-gonic/gin"
)

func JWTAuthHandler(c *gin.Context) {
	//获取请求头中的授权字段
	auth := c.Request.Header.Get("Authorization")
	if auth == "" {
		log.Println("授权字段为空")
		str.Response(c, str.ErrNotLogin, nil)
		c.Abort()
		return
	}
	//按空格分割取token
	token := strings.Split(auth, " ")[1]
	//解析token
	claims, err := jwt.ParseToken(token)
	if err != nil {
		log.Println("无效的token", err)
		str.Response(c, str.ErrInvalidAccessToken, nil)
		c.Abort()
		return
	}
	//将获取的用户id和用户名保存下来
	c.Set("userId", claims.UserID)
	c.Next()
}

func AdminAuthHandler(c *gin.Context) {
	//获取请求头中的授权字段
	token := c.Request.Header.Get("adminToken")
	if token == "" {
		log.Println("授权字段为空")
		str.Response(c, str.ErrNotLogin, nil)
		c.Abort()
		return
	}
	//解析token
	claims, err := jwt.ParseToken(token)
	if err != nil {
		log.Println("无效的token", err)
		str.Response(c, str.ErrNotLogin, nil)
		c.Abort()
		return
	}
	if claims.UserID != settings.Conf.Id {
		log.Println("错误的admin token")
		str.Response(c, str.ErrNotLogin, nil)
		c.Abort()
		return
	}
	c.Next()
}
