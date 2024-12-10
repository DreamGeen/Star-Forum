package httpHandler

import (
	"star/app/constant/str"
	"star/app/extra/tracing"
	"star/app/gateway/client"
	"star/app/utils/logging"
	"star/app/utils/request"
	"star/proto/message/messagePb"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func ListMessageCountHandler(c *gin.Context) {
	_, span := tracing.Tracer.Start(c.Request.Context(), "ListMessageCountHandler")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger := logging.LogServiceWithTrace(span, "GateWay.ListMessageCount")

	userId, err := request.GetUserId(c)
	if err != nil {
		logger.Error("get user id error",
			zap.Error(err))
		str.Response(c, err, nil)
		return
	}
	resp, err := client.ListMessageCount(c.Request.Context(),
		&messagePb.ListMessageCountRequest{
			UserId: userId,
		})
	if err != nil {
		logger.Error("list message count service error",
			zap.Error(err))
		str.Response(c, err, nil)
		return
	}

	str.Response(c, nil, map[string]interface{}{
		"count": resp.Count,
	})
}

func SendSystemMessageHandler(c *gin.Context) {
	_, span := tracing.Tracer.Start(c.Request.Context(), "SendSystemMessageHandler")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger := logging.LogServiceWithTrace(span, "GateWay.SendSystemMessage")

	req := &messagePb.SendSystemMessageRequest{
		ManagerId:   1820019310731464704,
		RecipientId: 0,
		Type:        "all",
		Title:       "Test",
		Content:     "test send system message",
	}
	_, err := client.SendSystemMessage(c.Request.Context(), req)
	if err != nil {
		logger.Error("send system message error",
			zap.Error(err))
		str.Response(c, err, nil)
		return
	}
	str.Response(c, nil, nil)
}

func SendPrivateMessageHandler(c *gin.Context) {
	_, span := tracing.Tracer.Start(c.Request.Context(), "SendPrivateMessageHandler")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger := logging.LogServiceWithTrace(span, "GateWay.SendPrivateMessage")

	recipientIdStr := c.Param("userId")
	recipientId, err := strconv.ParseInt(recipientIdStr, 10, 64)
	if err != nil {
		logger.Error("parse recipient id failed",
			zap.Error(err),
			zap.String("senderId", recipientIdStr))
		str.Response(c, str.ErrMessageError, nil)
		return
	}
	content := c.Query("content")
	senderId, err := request.GetUserId(c)
	if err != nil {
		logging.Logger.Error("get client id failed", zap.Error(err))
		str.Response(c, err, nil)
		return
	}
	_, err = client.SendPrivateMessage(c.Request.Context(), &messagePb.SendPrivateMessageRequest{
		SenderId:    senderId,
		RecipientId: recipientId,
		Content:     content,
	})
	if err != nil {
		logger.Error("send private message failed", zap.Error(err),
			zap.Int64("senderId", senderId),
			zap.Int64("recipientId", recipientId),
			zap.String("content", content))
		str.Response(c, err, nil)
		return
	}
	str.Response(c, nil, nil)
}
