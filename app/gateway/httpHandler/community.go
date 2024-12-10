package httpHandler

import (
	"star/app/constant/str"
	"star/app/extra/tracing"
	"star/app/gateway/client"
	"star/app/models"
	"star/app/utils/logging"
	"star/app/utils/request"
	"star/proto/community/communityPb"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func CreateCommunityHandler(c *gin.Context) {
	_, span := tracing.Tracer.Start(c.Request.Context(), "CreateCommunityHandler")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger := logging.LogServiceWithTrace(span, "GateWay.CreateCommunity")

	community := new(models.Community)
	if err := c.ShouldBindJSON(community); err != nil {
		logger.Error("invalid param",
			zap.Error(err))
		str.Response(c, str.ErrInvalidParam, nil)
		return
	}
	if err := validateDescriptionAndName(community); err != nil {
		logger.Error("invalid description",
			zap.Error(err))
		str.Response(c, err, nil)
		return
	}
	userId, err := request.GetUserId(c)
	if err != nil {
		logger.Error("user not log in",
			zap.Error(err))
		str.Response(c, err, nil)
		return
	}
	req := &communityPb.CreateCommunityRequest{
		CommunityName: community.CommunityName,
		Description:   community.Description,
		LeaderId:      userId,
	}
	if _, err := client.CreateCommunity(c.Request.Context(), req); err != nil {
		str.Response(c, err, nil)
		return
	}
	str.Response(c, nil, nil)
}

func GetFollowCommunityListHandler(c *gin.Context) {
	_, span := tracing.Tracer.Start(c.Request.Context(), "GetFollowCommunityListHandler")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger := logging.LogServiceWithTrace(span, "GateWay.GetFollowCommunityList")

	userId, err := request.GetUserId(c)
	if err != nil {
		logger.Warn("user not log in,but want to get follow community",
			zap.Error(err))
		str.Response(c, err, nil)
		return
	}
	resp, err := client.GetFollowCommunityList(c.Request.Context(), &communityPb.GetFollowCommunityListRequest{
		UserId: userId,
	})
	if err != nil {
		logger.Error("get follow community list service error",
			zap.Error(err),
			zap.Int64("userId", userId))
		str.Response(c, err, nil)
		return
	}

	str.Response(c, nil, map[string]interface{}{
		"communityList": resp.CommunityList,
	})
}
func FollowCommunityHandler(c *gin.Context) {
	_, span := tracing.Tracer.Start(c.Request.Context(), "FollowCommunityHandler")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger := logging.LogServiceWithTrace(span, "GateWay.FollowCommunityList")

	communityIdStr := c.Param("id")
	communityId, err := strconv.ParseInt(communityIdStr, 64, 10)
	if err != nil || communityId == 0 {
		logger.Error("follow user error,invalid param",
			zap.Error(err),
			zap.String("communityIdStr", communityIdStr))
		str.Response(c, str.ErrInvalidParam, nil)
		return
	}
	userId, err := request.GetUserId(c)
	if err != nil {
		logger.Warn("user not log in,but want to follow community",
			zap.Error(err),
			zap.Int64("communityId", communityId))
		str.Response(c, err, nil)
		return
	}
	_, err = client.FollowCommunity(c.Request.Context(), &communityPb.FollowCommunityRequest{
		ActorId:     userId,
		CommunityId: communityId,
	})
	if err != nil {
		logger.Error("follow community service error",
			zap.Error(err),
			zap.Int64("userId", userId),
			zap.Int64("communityId", communityId))
		str.Response(c, err, nil)
		return
	}
	str.Response(c, nil, nil)
	return
}

func UnFollowCommunityHandler(c *gin.Context) {
	_, span := tracing.Tracer.Start(c.Request.Context(), "UnFollowCommunityHandler")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger := logging.LogServiceWithTrace(span, "GateWay.UnFollowCommunity")

	communityIdStr := c.Query("id")
	communityId, err := strconv.ParseInt(communityIdStr, 64, 10)
	if err != nil || communityId == 0 {
		logger.Error("unfollow community error,invalid param",
			zap.Error(err),
			zap.String("communityIdStr", communityIdStr))
		str.Response(c, str.ErrInvalidParam, nil)
		return
	}
	userId, err := request.GetUserId(c)
	if err != nil {
		logger.Warn("user not log in,but want to unfollow community",
			zap.Error(err),
			zap.Int64("communityId", communityId))
		str.Response(c, err, nil)
		return
	}
	_, err = client.UnFollowCommunity(c.Request.Context(), &communityPb.UnFollowCommunityRequest{
		ActorId:     userId,
		CommunityId: communityId,
	})
	if err != nil {
		logging.Logger.Error("unfollow community service error",
			zap.Error(err),
			zap.Int64("userId", userId),
			zap.Int64("communityId", communityId))
		str.Response(c, err, nil)
		return
	}
	str.Response(c, nil, nil)
	return
}

func GetCommunityInfoHandler(c *gin.Context) {
	_, span := tracing.Tracer.Start(c.Request.Context(), "GetCommunityInfoHandler")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger := logging.LogServiceWithTrace(span, "GateWay.GetCommunityInfo")

	communityIdStr := c.Param("id")
	communityId, err := strconv.ParseInt(communityIdStr, 64, 10)
	if err != nil || communityId == 0 {
		logger.Error("invalid param",
			zap.Error(err),
			zap.String("communityIdStr", communityIdStr))
		str.Response(c, str.ErrInvalidParam, nil)
		return
	}
	resp, err := client.GetCommunityInfo(c.Request.Context(), &communityPb.GetCommunityInfoRequest{
		CommunityId: communityId,
	})
	if err != nil {
		logger.Error("get community info service error",
			zap.Error(err))
		str.Response(c, err, nil)
	}

	str.Response(c, nil, map[string]interface{}{
		"communityInfo": resp.Community,
	})
}

// validateDescription 效验简介字数
func validateDescriptionAndName(community *models.Community) error {
	if len(community.Description) < 2 {
		return str.ErrDescriptionShort
	}
	if len(community.Description) > 50 {
		return str.ErrDescriptionLong
	}
	if len(community.CommunityName) < 1 {
		return str.ErrCommunityNameEmpty
	}
	if len(community.CommunityName) > 10 {
		return str.ErrCommunityNameLong
	}
	return nil
}
