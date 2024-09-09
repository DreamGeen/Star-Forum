package web

import "github.com/gin-gonic/gin"

func MessageHandler(c *gin.Context) {
	c.Param("userId")

}
