package httpHandler

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	str2 "star/app/constant/str"
	"star/app/gateway/client"
	"star/app/gateway/models"
	"star/app/utils/logging"
	utils2 "star/app/utils/request"
	"star/proto/feed/feedPb"
)

func CreatePostHandler(c *gin.Context) {
	userId, err := utils2.GetUserId(c)
	if err != nil {
		logging.Logger.Warn("user not log in,but want to create feed")
		str2.Response(c, err, str2.Empty, nil)
		return
	}
	p := new(models.CreatePost)
	err = c.ShouldBindJSON(p)
	if err != nil {
		logging.Logger.Error("create feed error,invalid param",
			zap.Error(err),
			zap.Int64("userId", userId))
		str2.Response(c, str2.ErrInvalidParam, str2.Empty, nil)
		return
	}
	_, err = client.CreatePost(c, &feedPb.CreatePostRequest{
		UserId:      userId,
		IsScan:      p.IsScan,
		Content:     p.Content,
		CommunityId: p.CommunityId,
	})
	if err != nil {
		logging.Logger.Error("create feed  service error",
			zap.Error(err),
			zap.Int64("userId", userId),
			zap.Int64("communityId", p.CommunityId),
			zap.String("content", p.Content),
			zap.Bool("isScan", p.IsScan))
		str2.Response(c, err, str2.Empty, nil)
		return
	}
	str2.Response(c, nil, str2.Empty, nil)
	return
}

func GetPostByPopularityHandler(c *gin.Context) {
	p := new(models.GetCommunityPost)
	err := c.ShouldBindJSON(p)
	if err != nil {
		logging.Logger.Error("get feed by popularity error,invalid param",
			zap.Error(err),
		)
		str2.Response(c, str2.ErrInvalidParam, str2.Empty, nil)
		return
	}
	resp, err := client.GetPostByPopularity(c, &feedPb.GetPostListByPopularityRequest{
		Limit:       p.Limit,
		Page:        p.Page,
		CommunityId: p.CommunityId,
	})
	if err != nil {
		logging.Logger.Error("get feed by popularity service error",
			zap.Error(err),
			zap.Int64("limit", p.Limit),
			zap.Int64("page", p.Page),
			zap.Int64("communityId", p.CommunityId))
		str2.Response(c, err, str2.Empty, nil)
		return
	}
	str2.Response(c, nil, "posts", resp.Posts)
	return
}

func GetPostByTimeHandler(c *gin.Context) {
	p := new(models.GetCommunityPost)
	err := c.ShouldBindJSON(p)
	if err != nil {
		logging.Logger.Error("get feed by time error,invalid param",
			zap.Error(err),
		)
		str2.Response(c, str2.ErrInvalidParam, str2.Empty, nil)
		return
	}
	resp, err := client.GetPostByTime(c, &feedPb.GetPostListByTimeRequest{
		Limit:       p.Limit,
		Page:        p.Page,
		CommunityId: p.CommunityId,
	})
	if err != nil {
		logging.Logger.Error("get feed by time service error",
			zap.Error(err),
			zap.Int64("limit", p.Limit),
			zap.Int64("page", p.Page),
			zap.Int64("communityId", p.CommunityId))
		str2.Response(c, err, str2.Empty, nil)
		return
	}
	str2.Response(c, nil, "posts", resp.Posts)
	return
}
