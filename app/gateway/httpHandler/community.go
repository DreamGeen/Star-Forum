package httpHandler

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"log"
	"star/app/gateway/client"
	"star/constant/str"
	"star/models"
	"star/proto/community/communityPb"
	"star/utils"
	"strconv"
)

func CreateCommunityHandler(c *gin.Context) {
	community := new(models.Community)
	if err := c.ShouldBindJSON(community); err != nil {
		log.Println(err)
		str.Response(c, str.ErrInvalidParam, str.Empty, nil)
		return
	}
	if err := validateDescriptionAndName(community); err != nil {
		log.Println(err)
		str.Response(c, err, str.Empty, nil)
		return
	}
	userId, err := utils.GetUserId(c)
	if err != nil {
		log.Println(err)
		str.Response(c, err, str.Empty, nil)
		return
	}
	req := &communityPb.CreateCommunityRequest{
		CommunityName: community.CommunityName,
		Description:   community.Description,
		LeaderId:      userId,
	}
	if _, err := client.CreateCommunity(c, req); err != nil {
		str.Response(c, err, str.Empty, nil)
		return
	}
	str.Response(c, nil, str.Empty, nil)
}
func GetFollowCommunityList(c *gin.Context) {
	userId, err := utils.GetUserId(c)
	if err != nil {
		utils.Logger.Warn("user not log in,but want to get follow community",
			zap.Error(err))
		str.Response(c, err, str.Empty, nil)
		return
	}
	resp, err := client.GetFollowCommunityList(c, &communityPb.GetFollowCommunityListRequest{
		UserId: userId,
	})
	if err != nil {
		utils.Logger.Error("get follow community list service error",
			zap.Error(err),
			zap.Int64("userId", userId))
		str.Response(c, err, str.Empty, nil)
		return
	}
	str.Response(c, nil, "communityList", resp.CommunityList)
}
func FollowCommunity(c *gin.Context) {
	communityIdStr := c.Query("communityId")
	communityId, err := strconv.ParseInt(communityIdStr, 64, 10)
	if err != nil || communityId == 0 {
		utils.Logger.Error("follow user error,invalid param",
			zap.Error(err),
			zap.String("communityIdStr", communityIdStr))
		str.Response(c, str.ErrInvalidParam, str.Empty, nil)
		return
	}
	userId, err := utils.GetUserId(c)
	if err != nil {
		utils.Logger.Warn("user not log in,but want to follow community",
			zap.Error(err),
			zap.Int64("communityId", communityId))
		str.Response(c, err, str.Empty, nil)
		return
	}
	_, err = client.FollowCommunity(c, &communityPb.FollowCommunityRequest{
		ActorId:     userId,
		CommunityId: communityId,
	})
	if err != nil {
		utils.Logger.Error("follow community service error",
			zap.Error(err),
			zap.Int64("userId", userId),
			zap.Int64("communityId", communityId))
		str.Response(c, err, str.Empty, nil)
		return
	}
	str.Response(c, nil, str.Empty, nil)
	return
}

func UnFollowCommunity(c *gin.Context) {
	communityIdStr := c.Query("beFollowId")
	communityId, err := strconv.ParseInt(communityIdStr, 64, 10)
	if err != nil || communityId == 0 {
		utils.Logger.Error("unfollow community error,invalid param",
			zap.Error(err),
			zap.String("communityIdStr", communityIdStr))
		str.Response(c, str.ErrInvalidParam, str.Empty, nil)
		return
	}
	userId, err := utils.GetUserId(c)
	if err != nil {
		utils.Logger.Warn("user not log in,but want to unfollow community",
			zap.Error(err),
			zap.Int64("communityId", communityId))
		str.Response(c, err, str.Empty, nil)
		return
	}
	_, err = client.UnFollowCommunity(c, &communityPb.UnFollowCommunityRequest{
		ActorId:     userId,
		CommunityId: communityId,
	})
	if err != nil {
		utils.Logger.Error("unfollow community service error",
			zap.Error(err),
			zap.Int64("userId", userId),
			zap.Int64("communityId", communityId))
		str.Response(c, err, str.Empty, nil)
		return
	}
	str.Response(c, nil, str.Empty, nil)
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
