package httpHandler

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"star/app/gateway/client"
	"star/app/gateway/models"
	"star/constant/str"
	"star/proto/post/postPb"
	"star/utils"
)

func CreatePostHandler(c *gin.Context) {
	userId, err := utils.GetUserId(c)
	if err != nil {
		utils.Logger.Warn("user not log in,but want to create post")
		str.Response(c, err, str.Empty, nil)
		return
	}
	p := new(models.CreatePost)
	err = c.ShouldBindJSON(p)
	if err != nil {
		utils.Logger.Error("create post error,invalid param",
			zap.Error(err),
			zap.Int64("userId", userId))
		str.Response(c, str.ErrInvalidParam, str.Empty, nil)
		return
	}
	_, err = client.CreatePost(c, &postPb.CreatePostRequest{
		UserId:      userId,
		IsScan:      p.IsScan,
		Content:     p.Content,
		CommunityId: p.CommunityId,
	})
	if err != nil {
		utils.Logger.Error("create post  service error",
			zap.Error(err),
			zap.Int64("userId", userId),
			zap.Int64("communityId", p.CommunityId),
			zap.String("content", p.Content),
			zap.Bool("isScan", p.IsScan))
		str.Response(c, err, str.Empty, nil)
		return
	}
	str.Response(c, nil, str.Empty, nil)
	return
}

func GetPostByPopularityHandler(c *gin.Context) {
	p := new(models.GetCommunityPost)
	err := c.ShouldBindJSON(p)
	if err != nil {
		utils.Logger.Error("get post by popularity error,invalid param",
			zap.Error(err),
		)
		str.Response(c, str.ErrInvalidParam, str.Empty, nil)
		return
	}
	resp, err := client.GetPostByPopularity(c, &postPb.GetPostListByPopularityRequest{
		Limit:       p.Limit,
		Page:        p.Page,
		CommunityId: p.CommunityId,
	})
	if err != nil {
		utils.Logger.Error("get post by popularity service error",
			zap.Error(err),
			zap.Int64("limit", p.Limit),
			zap.Int64("page", p.Page),
			zap.Int64("communityId", p.CommunityId))
		str.Response(c, err, str.Empty, nil)
		return
	}
	str.Response(c, nil, "posts", resp.Posts)
	return
}

func GetPostByTimeHandler(c *gin.Context) {
	p := new(models.GetCommunityPost)
	err := c.ShouldBindJSON(p)
	if err != nil {
		utils.Logger.Error("get post by time error,invalid param",
			zap.Error(err),
		)
		str.Response(c, str.ErrInvalidParam, str.Empty, nil)
		return
	}
	resp, err := client.GetPostByTime(c, &postPb.GetPostListByTimeRequest{
		Limit:       p.Limit,
		Page:        p.Page,
		CommunityId: p.CommunityId,
	})
	if err != nil {
		utils.Logger.Error("get post by time service error",
			zap.Error(err),
			zap.Int64("limit", p.Limit),
			zap.Int64("page", p.Page),
			zap.Int64("communityId", p.CommunityId))
		str.Response(c, err, str.Empty, nil)
		return
	}
	str.Response(c, nil, "posts", resp.Posts)
	return
}
