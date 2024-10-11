package httpHandler

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"log"
	str2 "star/app/constant/str"
	"star/app/gateway/client"
	"star/app/models"
	"star/app/utils/logging"
	utils2 "star/app/utils/request"
	"star/proto/community/communityPb"
	"strconv"
)

func CreateCommunityHandler(c *gin.Context) {
	community := new(models.Community)
	if err := c.ShouldBindJSON(community); err != nil {
		log.Println(err)
		str2.Response(c, str2.ErrInvalidParam, str2.Empty, nil)
		return
	}
	if err := validateDescriptionAndName(community); err != nil {
		log.Println(err)
		str2.Response(c, err, str2.Empty, nil)
		return
	}
	userId, err := utils2.GetUserId(c)
	if err != nil {
		log.Println(err)
		str2.Response(c, err, str2.Empty, nil)
		return
	}
	req := &communityPb.CreateCommunityRequest{
		CommunityName: community.CommunityName,
		Description:   community.Description,
		LeaderId:      userId,
	}
	if _, err := client.CreateCommunity(c, req); err != nil {
		str2.Response(c, err, str2.Empty, nil)
		return
	}
	str2.Response(c, nil, str2.Empty, nil)
}
func GetFollowCommunityList(c *gin.Context) {
	userId, err := utils2.GetUserId(c)
	if err != nil {
		logging.Logger.Warn("user not log in,but want to get follow community",
			zap.Error(err))
		str2.Response(c, err, str2.Empty, nil)
		return
	}
	resp, err := client.GetFollowCommunityList(c, &communityPb.GetFollowCommunityListRequest{
		UserId: userId,
	})
	if err != nil {
		logging.Logger.Error("get follow community list service error",
			zap.Error(err),
			zap.Int64("userId", userId))
		str2.Response(c, err, str2.Empty, nil)
		return
	}
	str2.Response(c, nil, "communityList", resp.CommunityList)
}
func FollowCommunity(c *gin.Context) {
	communityIdStr := c.Query("communityId")
	communityId, err := strconv.ParseInt(communityIdStr, 64, 10)
	if err != nil || communityId == 0 {
		logging.Logger.Error("follow user error,invalid param",
			zap.Error(err),
			zap.String("communityIdStr", communityIdStr))
		str2.Response(c, str2.ErrInvalidParam, str2.Empty, nil)
		return
	}
	userId, err := utils2.GetUserId(c)
	if err != nil {
		logging.Logger.Warn("user not log in,but want to follow community",
			zap.Error(err),
			zap.Int64("communityId", communityId))
		str2.Response(c, err, str2.Empty, nil)
		return
	}
	_, err = client.FollowCommunity(c, &communityPb.FollowCommunityRequest{
		ActorId:     userId,
		CommunityId: communityId,
	})
	if err != nil {
		logging.Logger.Error("follow community service error",
			zap.Error(err),
			zap.Int64("userId", userId),
			zap.Int64("communityId", communityId))
		str2.Response(c, err, str2.Empty, nil)
		return
	}
	str2.Response(c, nil, str2.Empty, nil)
	return
}

func UnFollowCommunity(c *gin.Context) {
	communityIdStr := c.Query("beFollowId")
	communityId, err := strconv.ParseInt(communityIdStr, 64, 10)
	if err != nil || communityId == 0 {
		logging.Logger.Error("unfollow community error,invalid param",
			zap.Error(err),
			zap.String("communityIdStr", communityIdStr))
		str2.Response(c, str2.ErrInvalidParam, str2.Empty, nil)
		return
	}
	userId, err := utils2.GetUserId(c)
	if err != nil {
		logging.Logger.Warn("user not log in,but want to unfollow community",
			zap.Error(err),
			zap.Int64("communityId", communityId))
		str2.Response(c, err, str2.Empty, nil)
		return
	}
	_, err = client.UnFollowCommunity(c, &communityPb.UnFollowCommunityRequest{
		ActorId:     userId,
		CommunityId: communityId,
	})
	if err != nil {
		logging.Logger.Error("unfollow community service error",
			zap.Error(err),
			zap.Int64("userId", userId),
			zap.Int64("communityId", communityId))
		str2.Response(c, err, str2.Empty, nil)
		return
	}
	str2.Response(c, nil, str2.Empty, nil)
	return
}

//func GetCommunityListHandler(c *gin.Context) {
//	resp, err := client.GetCommunityList()
//	if err != nil {
//		str.Response(c, err, str.Empty, nil)
//		return
//	}
//	str.Response(c, nil, "communitys", resp.Communitys)
//}

//func ShowCommunityHandler(c *gin.Context) {
//	communityIdString := c.Param("id")
//	communityId, _ := strconv.Atoi(communityIdString)
//	resp, err := client.ShowCommunity()
//	if err != nil {
//		str.Response(c, err, str.Empty, nil)
//		return
//	}
//	str.Response(c, nil, str.Empty, resp)
//}

// validateDescription 效验简介字数
func validateDescriptionAndName(community *models.Community) error {
	if len(community.Description) < 2 {
		return str2.ErrDescriptionShort
	}
	if len(community.Description) > 50 {
		return str2.ErrDescriptionLong
	}
	if len(community.CommunityName) < 1 {
		return str2.ErrCommunityNameEmpty
	}
	if len(community.CommunityName) > 10 {
		return str2.ErrCommunityNameLong
	}
	return nil
}
