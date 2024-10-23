package httpHandler

import (
	"star/app/constant/str"
	"star/app/extra/tracing"
	"star/app/gateway/client"
	"star/app/utils/logging"
	"star/app/utils/request"
	"star/proto/relation/relationPb"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func FollowHandler(c *gin.Context) {
	_,span:=tracing.Tracer.Start(c.Request.Context(),"FollowHandler")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger:=logging.LogServiceWithTrace(span,"GateWay.Follow")

	beFollowIdStr := c.Query("beFollowId")
	beFollowId, err := strconv.ParseInt(beFollowIdStr, 64, 10)
	if err != nil || beFollowId == 0 {
		logger.Error("follow user error,invalid param",
			zap.Error(err),
			zap.String("beFollowIdStr", beFollowIdStr))
		str.Response(c, str.ErrInvalidParam, str.Empty, nil)
		return
	}
	userId, err := request.GetUserId(c)
	if err != nil {
		logger.Warn("user not log in,but want to follow",
			zap.Error(err),
			zap.Int64("beFollowId", beFollowId))
		str.Response(c, err, str.Empty, nil)
		return
	}
	_, err = client.Follow(c.Request.Context(), &relationPb.FollowRequest{
		UserId:       userId,
		BeFollowerId: beFollowId,
	})
	if err != nil {
		logger.Error("follow service error",
			zap.Error(err),
			zap.Int64("userId", userId),
			zap.Int64("beFollowId", beFollowId))
		str.Response(c, err, str.Empty, nil)
		return
	}
	str.Response(c, nil, str.Empty, nil)
	return
}

func UnFollowHandler(c *gin.Context) {
	_,span:=tracing.Tracer.Start(c.Request.Context(),"UnFollowHandler")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger:=logging.LogServiceWithTrace(span,"GateWay.UnFollow")

	unBeFollowIdStr := c.Query("beFollowId")
	unBeFollowId, err := strconv.ParseInt(unBeFollowIdStr, 64, 10)
	if err != nil || unBeFollowId == 0 {
		logger.Error("unfollow user error,invalid param",
			zap.Error(err),
			zap.String("UnBeFollowIdStr", unBeFollowIdStr))
		str.Response(c, str.ErrInvalidParam, str.Empty, nil)
		return
	}
	userId, err := request.GetUserId(c)
	if err != nil {
		logger.Warn("user not log in,but want to unfollow",
			zap.Error(err),
			zap.Int64("unBeFollowId", unBeFollowId))
		str.Response(c, err, str.Empty, nil)
		return
	}
	_, err = client.UnFollow(c.Request.Context(), &relationPb.UnFollowRequest{
		UserId:         userId,
		UnBeFollowerId: unBeFollowId,
	})
	if err != nil {
		logger.Error("unfollow service error",
			zap.Error(err),
			zap.Int64("userId", userId),
			zap.Int64("unBeFollowId", unBeFollowId))
		str.Response(c, err, str.Empty, nil)
		return
	}
	str.Response(c, nil, str.Empty, nil)
	return
}

func GetFollowListHandler(c *gin.Context) {
	_,span:=tracing.Tracer.Start(c.Request.Context(),"GetFollowListHandler")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger:=logging.LogServiceWithTrace(span,"GateWay.GetFollowList")

	userId, err := request.GetUserId(c)
	if err != nil {
		logging.Logger.Warn("user not log in,but want to get follow list",
			zap.Error(err))
		str.Response(c, err, str.Empty, nil)
		return
	}
	resp, err := client.GetFollowList(c.Request.Context(), &relationPb.GetFollowListRequest{
		UserId: userId,
	})
	if err != nil {
		logger.Error("get follow list service error",
			zap.Error(err),
			zap.Int64("userId", userId))
		str.Response(c, err, str.Empty, nil)
		return
	}
	str.Response(c, nil, "followList", resp.FollowList)
	return
}

func GetFansListHandler(c *gin.Context) {
	_,span:=tracing.Tracer.Start(c.Request.Context(),"GetFansListHandler")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger:=logging.LogServiceWithTrace(span,"GateWay.GetFansList")

	userId, err := request.GetUserId(c)
	if err != nil {
		logger.Warn("user not log in,but want to get fans list",
			zap.Error(err))
		str.Response(c, err, str.Empty, nil)
		return
	}
	resp, err := client.GetFansList(c.Request.Context(), &relationPb.GetFansListRequest{
		UserId: userId,
	})
	if err != nil {
		logger.Error("get fans list service error",
			zap.Error(err),
			zap.Int64("userId", userId))
		str.Response(c, err, str.Empty, nil)
		return
	}
	str.Response(c, nil, "fansList", resp.FansList)
	return
}
