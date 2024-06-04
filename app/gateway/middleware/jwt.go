package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"

	"star/utils"
)

func JWTAuthHandler(c *gin.Context) {
	//获取请求头中的授权字段
	auth := c.Request.Header.Get("Authorization")
	if auth == "" {
		utils.ResponseMessage(c, utils.CodeNotLogin)
		c.Abort()
	}
	//按空格分割取token
	token := strings.Split(auth, " ")[1]
	//解析token
	claims, err := utils.ParseToken(token)
	if err != nil {
		utils.ResponseMessageWithData(c, utils.CodeNotLogin, err)
		c.Abort()
	}
	//将获取的用户id和用户名保存下来
	c.Set("userid", claims.UserID)
	c.Set("username", claims.UserName)
	c.Next()
}
