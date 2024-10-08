package httpHandler

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"star/app/gateway/client"
	"star/constant/str"
	"star/proto/relation/relationPb"
	"star/utils"
	"strconv"
)

func FollowHandler(c *gin.Context) {
	beFollowIdStr := c.Query("beFollowId")
	beFollowId, err := strconv.ParseInt(beFollowIdStr, 64, 10)
	if err != nil || beFollowId == 0 {
		utils.Logger.Error("follow user error,invalid param",
			zap.Error(err),
			zap.String("beFollowIdStr", beFollowIdStr))
		str.Response(c, str.ErrInvalidParam, str.Empty, nil)
		return
	}
	userId, err := utils.GetUserId(c)
	if err != nil {
		utils.Logger.Warn("user not log in,but want to follow",
			zap.Error(err),
			zap.Int64("beFollowId", beFollowId))
		str.Response(c, err, str.Empty, nil)
		return
	}
	_, err = client.Follow(c, &relationPb.FollowRequest{
		UserId:       userId,
		BeFollowerId: beFollowId,
	})
	if err != nil {
		utils.Logger.Error("follow service error",
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
	unBeFollowIdStr := c.Query("beFollowId")
	unBeFollowId, err := strconv.ParseInt(unBeFollowIdStr, 64, 10)
	if err != nil || unBeFollowId == 0 {
		utils.Logger.Error("unfollow user error,invalid param",
			zap.Error(err),
			zap.String("UnBeFollowIdStr", unBeFollowIdStr))
		str.Response(c, str.ErrInvalidParam, str.Empty, nil)
		return
	}
	userId, err := utils.GetUserId(c)
	if err != nil {
		utils.Logger.Warn("user not log in,but want to unfollow",
			zap.Error(err),
			zap.Int64("unBeFollowId", unBeFollowId))
		str.Response(c, err, str.Empty, nil)
		return
	}
	_, err = client.UnFollow(c, &relationPb.UnFollowRequest{
		UserId:         userId,
		UnBeFollowerId: unBeFollowId,
	})
	if err != nil {
		utils.Logger.Error("unfollow service error",
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
	userId, err := utils.GetUserId(c)
	if err != nil {
		utils.Logger.Warn("user not log in,but want to get follow list",
			zap.Error(err))
		str.Response(c, err, str.Empty, nil)
		return
	}
	resp, err := client.GetFollowList(c, &relationPb.GetFollowRequest{
		UserId: userId,
	})
	if err != nil {
		utils.Logger.Error("get follow list service error",
			zap.Error(err),
			zap.Int64("userId", userId))
		str.Response(c, err, str.Empty, nil)
		return
	}
	str.Response(c, nil, "followList", resp.FollowList)
	return
}

func GetFansListHandler(c *gin.Context) {
	userId, err := utils.GetUserId(c)
	if err != nil {
		utils.Logger.Warn("user not log in,but want to get fans list",
			zap.Error(err))
		str.Response(c, err, str.Empty, nil)
		return
	}
	resp, err := client.GetFansList(c, &relationPb.GetFansListRequest{
		UserId: userId,
	})
	if err != nil {
		utils.Logger.Error("get fans list service error",
			zap.Error(err),
			zap.Int64("userId", userId))
		str.Response(c, err, str.Empty, nil)
		return
	}
	str.Response(c, nil, "fansList", resp.FansList)
	return
}
