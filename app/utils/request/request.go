package request

import (
	"github.com/gin-gonic/gin"
	"star/app/constant/str"
)

func GetUserId(c *gin.Context) (int64, error) {
	userId, ok := c.Get("userId")
	if !ok {
		return 0, str.ErrNotLogin
	}
	return userId.(int64), nil
}
