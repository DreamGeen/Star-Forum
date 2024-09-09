package middleware

import (
	"log"
	"star/constant/str"
	"star/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

func JWTAuthHandler(c *gin.Context) {
	//获取请求头中的授权字段
	auth := c.Request.Header.Get("Authorization")
	if auth == "" {
		log.Println("授权字段为空")
		str.Response(c, str.ErrNotLogin, str.Empty, nil)
		c.Abort()
	}
	//按空格分割取token
	token := strings.Split(auth, " ")[1]
	//解析token
	claims, err := utils.ParseToken(token)
	if err != nil {
		log.Println("无效的token", err)
		str.Response(c, str.ErrNotLogin, str.Empty, nil)
		c.Abort()
	}
	//将获取的用户id和用户名保存下来
	c.Set("userId", claims.UserID)
	c.Next()
}
