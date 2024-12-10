package httpHandler

import (
	"star/app/constant/str"
	"star/app/extra/tracing"
	"star/app/gateway/client"
	"star/app/gateway/models"
	"star/app/utils/logging"
	"star/app/utils/request"
	"star/proto/collect/collectPb"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func CollectActionHandler(c *gin.Context) {
	_, span := tracing.Tracer.Start(c.Request.Context(), "CollectActionHandler")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger := logging.LogServiceWithTrace(span, "GateWay.CollectAction")

	co := new(models.CollectAction)
	err := c.ShouldBindJSON(co)
	if err != nil {
		logger.Error("collect action error,invalid param",
			zap.Error(err))
		str.Response(c, str.ErrInvalidParam, nil)
		return
	}
	co.ActionId, err = request.GetUserId(c)
	if err != nil {
		logger.Warn("user not log in,but want to collect action",
			zap.Error(err))
		str.Response(c, err, nil)
		return
	}
	_, err = client.CollectAction(c.Request.Context(), &collectPb.CollectActionRequest{
		ActorId:    co.ActionId,
		PostId:     co.PostId,
		ActionType: co.ActionType,
	})
	if err != nil {
		logging.Logger.Error("collect action service error",
			zap.Error(err),
			zap.Int64("userId", co.ActionId),
			zap.Int64("postId", co.PostId),
			zap.Uint32("actionType", co.ActionType))
		str.Response(c, err, nil)
		return
	}
	str.Response(c, nil, nil)
	return
}

func CollectListHandler(c *gin.Context) {
	_, span := tracing.Tracer.Start(c.Request.Context(), "CollectListHandler")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger := logging.LogServiceWithTrace(span, "GateWay.CollectList")

	userId, err := request.GetUserId(c)
	if err != nil {
		logger.Warn("user not log in,but want to list collect feed",
			zap.Error(err))
		str.Response(c, err, nil)
		return
	}
	resp, err := client.CollectList(c.Request.Context(), &collectPb.CollectListRequest{
		ActorId: userId,
	})
	if err != nil {
		logger.Error("collect list service error",
			zap.Error(err),
			zap.Int64("userId", userId))
		str.Response(c, err, nil)
		return
	}
	str.Response(c, nil, map[string]interface{}{
		"posts": resp.Posts,
	})
	return

}
