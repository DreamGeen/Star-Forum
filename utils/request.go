package utils

import (
	"github.com/gin-gonic/gin"
	"star/constant/str"
	"star/models"
)

func GetUser(c *gin.Context) (*models.User, error) {
	userId, ok := c.Get("userId")
	if !ok {
		return nil, str.ErrNotLogin
	}
	username, ok := c.Get("username")
	if !ok {
		return nil, str.ErrNotLogin
	}
	img, ok := c.Get("img")
	if !ok {
		return nil, str.ErrNotLogin
	}
	return &models.User{
		UserId:   userId.(int64),
		Username: username.(string),
		Img:      img.(string),
	}, nil
}
