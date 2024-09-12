package httpHandler

import (
	"github.com/gin-gonic/gin"
	websocket2 "star/app/channels/websocket"
	"star/constant/str"
	"star/models"
	"star/utils"
	"strconv"
)

func ChatCommunityHandler(c *gin.Context) {
	//获取communityId
	communityIdStr := c.Param("communityId")
	communityId, err := strconv.ParseInt(communityIdStr, 10, 64)
	if err != nil {
		str.Response(c, str.ErrInvalidParam, str.Empty, nil)
		return
	}
	//获取userId
	userId, err := utils.GetUserId(c)
	if err != nil {
		str.Response(c, str.ErrInvalidParam, str.Empty, nil)
		return
	}
	user := &models.User{UserId: userId}
	hubManager := websocket2.NewHubManager()
	websocket2.ServeWs(communityId, communityIdStr, user, hubManager, c.Writer, c.Request)
}
