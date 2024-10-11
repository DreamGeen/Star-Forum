package httpHandler

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	str2 "star/app/constant/str"
	"star/app/gateway/client"
	"star/app/utils/logging"
	utils2 "star/app/utils/request"
	"star/proto/relation/relationPb"
	"strconv"
)

func FollowHandler(c *gin.Context) {
	beFollowIdStr := c.Query("beFollowId")
	beFollowId, err := strconv.ParseInt(beFollowIdStr, 64, 10)
	if err != nil || beFollowId == 0 {
		logging.Logger.Error("follow user error,invalid param",
			zap.Error(err),
			zap.String("beFollowIdStr", beFollowIdStr))
		str2.Response(c, str2.ErrInvalidParam, str2.Empty, nil)
		return
	}
	userId, err := utils2.GetUserId(c)
	if err != nil {
		logging.Logger.Warn("user not log in,but want to follow",
			zap.Error(err),
			zap.Int64("beFollowId", beFollowId))
		str2.Response(c, err, str2.Empty, nil)
		return
	}
	_, err = client.Follow(c, &relationPb.FollowRequest{
		UserId:       userId,
		BeFollowerId: beFollowId,
	})
	if err != nil {
		logging.Logger.Error("follow service error",
			zap.Error(err),
			zap.Int64("userId", userId),
			zap.Int64("beFollowId", beFollowId))
		str2.Response(c, err, str2.Empty, nil)
		return
	}
	str2.Response(c, nil, str2.Empty, nil)
	return
}

func UnFollowHandler(c *gin.Context) {
	unBeFollowIdStr := c.Query("beFollowId")
	unBeFollowId, err := strconv.ParseInt(unBeFollowIdStr, 64, 10)
	if err != nil || unBeFollowId == 0 {
		logging.Logger.Error("unfollow user error,invalid param",
			zap.Error(err),
			zap.String("UnBeFollowIdStr", unBeFollowIdStr))
		str2.Response(c, str2.ErrInvalidParam, str2.Empty, nil)
		return
	}
	userId, err := utils2.GetUserId(c)
	if err != nil {
		logging.Logger.Warn("user not log in,but want to unfollow",
			zap.Error(err),
			zap.Int64("unBeFollowId", unBeFollowId))
		str2.Response(c, err, str2.Empty, nil)
		return
	}
	_, err = client.UnFollow(c, &relationPb.UnFollowRequest{
		UserId:         userId,
		UnBeFollowerId: unBeFollowId,
	})
	if err != nil {
		logging.Logger.Error("unfollow service error",
			zap.Error(err),
			zap.Int64("userId", userId),
			zap.Int64("unBeFollowId", unBeFollowId))
		str2.Response(c, err, str2.Empty, nil)
		return
	}
	str2.Response(c, nil, str2.Empty, nil)
	return
}

func GetFollowListHandler(c *gin.Context) {
	userId, err := utils2.GetUserId(c)
	if err != nil {
		logging.Logger.Warn("user not log in,but want to get follow list",
			zap.Error(err))
		str2.Response(c, err, str2.Empty, nil)
		return
	}
	resp, err := client.GetFollowList(c, &relationPb.GetFollowRequest{
		UserId: userId,
	})
	if err != nil {
		logging.Logger.Error("get follow list service error",
			zap.Error(err),
			zap.Int64("userId", userId))
		str2.Response(c, err, str2.Empty, nil)
		return
	}
	str2.Response(c, nil, "followList", resp.FollowList)
	return
}

func GetFansListHandler(c *gin.Context) {
	userId, err := utils2.GetUserId(c)
	if err != nil {
		logging.Logger.Warn("user not log in,but want to get fans list",
			zap.Error(err))
		str2.Response(c, err, str2.Empty, nil)
		return
	}
	resp, err := client.GetFansList(c, &relationPb.GetFansListRequest{
		UserId: userId,
	})
	if err != nil {
		logging.Logger.Error("get fans list service error",
			zap.Error(err),
			zap.Int64("userId", userId))
		str2.Response(c, err, str2.Empty, nil)
		return
	}
	str2.Response(c, nil, "fansList", resp.FansList)
	return
}
