package httpHandler

import (
	"star/app/constant/str"
	"star/app/extra/tracing"
	"star/app/gateway/client"
	"star/app/gateway/models"
	"star/app/utils/logging"
	"star/app/utils/request"
	"star/proto/feed/feedPb"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func GetCommunityPostByTimeHandler(c *gin.Context) {
	_, span := tracing.Tracer.Start(c.Request.Context(), "GetCommunityPostByTimeHandler")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger := logging.LogServiceWithTrace(span, "GateWay.GetCommunityPostByTime")

	p := new(models.GetCommunityPost)
	err := c.ShouldBindJSON(p)
	if err != nil {
		logger.Error("get feed by time error,invalid param",
			zap.Error(err),
		)
		str.Response(c, str.ErrInvalidParam, nil)
		return
	}
	resp, err := client.GetCommunityPostByTime(c.Request.Context(), &feedPb.GetCommunityPostByTimeRequest{
		CommunityId: p.CommunityId,
		Page:        p.Page,
		ActorId:     p.ActorId,
		LastPostId:  p.LastPostId,
	})
	if err != nil {
		logger.Error("get feed by time service error",
			zap.Error(err),
			zap.Int64("actorId", p.ActorId),
			zap.Int64("page", p.Page),
			zap.Int64("communityId", p.CommunityId),
			zap.Int64("lastPostId", p.LastPostId))
		str.Response(c, err, nil)
		return
	}

	str.Response(c, nil, map[string]interface{}{
		"posts": resp.Posts,
	})
	return
}

func GetCommunityPostByNewReplyHandler(c *gin.Context) {
	_, span := tracing.Tracer.Start(c.Request.Context(), "GetCommunityPostByNewReplyHandler")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger := logging.LogServiceWithTrace(span, "GateWay.GetCommunityPostByNewReply")

	p := new(models.GetCommunityPostByNewReply)
	err := c.ShouldBindJSON(p)
	if err != nil {
		logger.Error("get feed by time error,invalid param",
			zap.Error(err),
		)
		str.Response(c, str.ErrInvalidParam, nil)
		return
	}
	resp, err := client.GetCommunityPostByNewRely(c.Request.Context(), &feedPb.GetCommunityPostByNewReplyRequest{
		CommunityId:   p.CommunityId,
		ActorId:       p.ActorId,
		Page:          p.Page,
		LastReplyTime: p.LastRelyTime,
	})
	if err != nil {
		logger.Error("get feed by time service error",
			zap.Error(err),
			zap.Int64("actorId", p.ActorId),
			zap.Int64("page", p.Page),
			zap.Int64("communityId", p.CommunityId))
		str.Response(c, err, nil)
		return
	}

	str.Response(c, nil, map[string]interface{}{
		"posts": resp.Posts,
	})
	return
}

func GetPostByNewRelationHandler(c *gin.Context) {
	_, span := tracing.Tracer.Start(c.Request.Context(), "GetPostByNewRelationHandler")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger := logging.LogServiceWithTrace(span, "GateWay.GetPostByNewRelation")

	pageStr := c.Query("page")
	page, err := strconv.ParseInt(pageStr, 10, 64)
	if err != nil {
		logger.Error("invalid param",
			zap.Error(err),
			zap.String("pageStr", pageStr))
		str.Response(c, str.ErrInvalidParam, nil)
		return
	}
	userId, err := request.GetUserId(c)
	if err != nil {
		logger.Error("get userId error",
			zap.Error(err))
		str.Response(c, err, nil)
		return
	}
	resp, err := client.GetPostByRelation(c.Request.Context(), &feedPb.GetPostByRelationRequest{
		ActorId: userId,
		Page:    page,
	})
	if err != nil {
		logger.Error("get feed by time service error",
			zap.Error(err),
			zap.Int64("actorId", userId),
			zap.Int64("page", page))
		str.Response(c, err, nil)
		return
	}
	str.Response(c, nil, map[string]interface{}{
		"posts": resp.Posts,
	})
	return
}
