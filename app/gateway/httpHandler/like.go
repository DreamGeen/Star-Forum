package httpHandler

import (
	"star/app/constant/str"
	"star/app/extra/tracing"
	"star/app/gateway/client"
	"star/app/gateway/models"
	"star/app/utils/logging"
	"star/app/utils/request"
	"star/proto/like/likePb"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func LikeActionHandler(c *gin.Context) {
	_,span:=tracing.Tracer.Start(c.Request.Context(),"LikeActionHandler")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger:=logging.LogServiceWithTrace(span,"GateWay.LikeAction")

	l := new(models.LikeAction)
	err := c.ShouldBindJSON(l)
	if err != nil {
		logger.Error("like action error,invalid param",
			zap.Error(err))
		str.Response(c, str.ErrInvalidParam, str.Empty, nil)
		return
	}
	l.UserId, err = request.GetUserId(c)
	if err != nil {
		logger.Warn("user not log in,but want to like action",
			zap.Error(err))
		str.Response(c, err, str.Empty, nil)
		return
	}
	l.Url = c.Request.URL.RawPath
	_, err = client.LikeAction(c.Request.Context(), &likePb.LikeActionRequest{
		UserId:     l.UserId,
		SourceId:   l.SourceId,
		SourceType: l.SourceType,
		ActionTye:  l.ActionType,
		Url:        l.Url,
	})
	if err != nil {
		logger.Error("like action service error",
			zap.Error(err),
			zap.Int64("userId", l.UserId),
			zap.Int64("sourceId", l.SourceId),
			zap.Uint32("sourceType", l.SourceType),
			zap.String("url", l.Url))
		str.Response(c, err, str.Empty, nil)
		return
	}
	str.Response(c, nil, str.Empty, nil)
	return
}

func LikeListHandler(c *gin.Context) {
	_,span:=tracing.Tracer.Start(c.Request.Context(),"LikeListHandler")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger:=logging.LogServiceWithTrace(span,"GateWay.LikeList")

	userId, err := request.GetUserId(c)
	if err != nil {
		logger.Warn("user not log in,but want to list like feed",
			zap.Error(err))
		str.Response(c, err, str.Empty, nil)
		return
	}
	resp, err := client.LikeList(c.Request.Context(), &likePb.LikeListRequest{
		UserId: userId,
	})
	if err != nil {
		logger.Error("like list service error",
			zap.Error(err),
			zap.Int64("userId", userId))
		str.Response(c, err, str.Empty, nil)
		return
	}
	str.Response(c, nil, "posts", resp.Posts)
	return
}
