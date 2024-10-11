package httpHandler

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	str2 "star/app/constant/str"
	"star/app/gateway/client"
	"star/app/utils/logging"
	utils2 "star/app/utils/request"
	"star/proto/message/messagePb"
	"strconv"
)

func ListMessageCountHandler(c *gin.Context) {
	userId, err := utils2.GetUserId(c)
	if err != nil {
		logging.Logger.Error("get user id error", zap.Error(err))
		str2.Response(c, err, str2.Empty, nil)
		return
	}
	resp, err := client.ListMessageCount(c,
		&messagePb.ListMessageCountRequest{
			UserId: userId,
		})
	if err != nil {
		str2.Response(c, err, str2.Empty, nil)
		return
	}
	str2.Response(c, nil, "count", resp.Count)
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
		logging.Logger.Error("send system message error", zap.Error(err))
		str2.Response(c, err, str2.Empty, nil)
		return
	}
	str2.Response(c, nil, "success", resp)
}

func SendPrivateMessageHandler(c *gin.Context) {
	recipientIdStr := c.Param("userId")
	recipientId, err := strconv.ParseInt(recipientIdStr, 10, 64)
	if err != nil {
		logging.Logger.Error("parse recipient id failed", zap.Error(err), zap.String("senderId", recipientIdStr))
		str2.Response(c, str2.ErrMessageError, str2.Empty, nil)
		return
	}
	content := c.Query("content")
	senderId, err := utils2.GetUserId(c)
	if err != nil {
		logging.Logger.Error("get client id failed", zap.Error(err))
		str2.Response(c, err, str2.Empty, nil)
		return
	}
	_, err = client.SendPrivateMessage(c, &messagePb.SendPrivateMessageRequest{
		SenderId:    senderId,
		RecipientId: recipientId,
		Content:     content,
	})
	if err != nil {
		logging.Logger.Error("send private message failed", zap.Error(err),
			zap.Int64("senderId", senderId),
			zap.Int64("recipientId", recipientId),
			zap.String("content", content))
		str2.Response(c, err, str2.Empty, nil)
		return
	}
	str2.Response(c, nil, str2.Empty, nil)
}
