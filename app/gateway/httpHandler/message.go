package httpHandler

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"star/app/gateway/client"
	"star/constant/str"
	"star/proto/message/messagePb"
	"star/utils"
	"strconv"
)

func ListMessageCountHandler(c *gin.Context) {
	userId, err := utils.GetUserId(c)
	if err != nil {
		utils.Logger.Error("get user id error", zap.Error(err))
		str.Response(c, err, str.Empty, nil)
		return
	}
	resp, err := client.ListMessageCount(c,
		&messagePb.ListMessageCountRequest{
			UserId: userId,
		})
	if err != nil {
		str.Response(c, err, str.Empty, nil)
		return
	}
	str.Response(c, nil, "count", resp.Count)
}

func SendSystemHandler(c *gin.Context) {
	req := &messagePb.SendSystemMessageRequest{
		ManagerId:   1820019310731464704,
		RecipientId: 0,
		Type:        "all",
		Title:       "Test",
		Content:     "test send system message",
	}
	resp, err := client.SendSystemMessage(c, req)
	if err != nil {
		utils.Logger.Error("send system message error", zap.Error(err))
		str.Response(c, err, str.Empty, nil)
		return
	}
	str.Response(c, nil, "success", resp)
}

func SendPrivateMessageHandler(c *gin.Context) {
	recipientIdStr := c.Param("userId")
	recipientId, err := strconv.ParseInt(recipientIdStr, 10, 64)
	if err != nil {
		utils.Logger.Error("parse recipient id failed", zap.Error(err), zap.String("senderId", recipientIdStr))
		str.Response(c, str.ErrMessageError, str.Empty, nil)
		return
	}
	content := c.Query("content")
	senderId, err := utils.GetUserId(c)
	if err != nil {
		utils.Logger.Error("get client id failed", zap.Error(err))
		str.Response(c, err, str.Empty, nil)
		return
	}
	_, err = client.SendPrivateMessage(c, &messagePb.SendPrivateMessageRequest{
		SenderId:    senderId,
		RecipientId: recipientId,
		Content:     content,
	})
	if err != nil {
		utils.Logger.Error("send private message failed", zap.Error(err),
			zap.Int64("senderId", senderId),
			zap.Int64("recipientId", recipientId),
			zap.String("content", content))
		str.Response(c, err, str.Empty, nil)
		return
	}
	str.Response(c, nil, str.Empty, nil)
}
