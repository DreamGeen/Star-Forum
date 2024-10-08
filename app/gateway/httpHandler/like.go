package httpHandler

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"star/app/gateway/client"
	"star/app/gateway/models"
	"star/constant/str"
	"star/proto/like/likePb"
	"star/utils"
)

func LikeActionHandler(c *gin.Context) {
	l := new(models.LikeAction)
	err := c.ShouldBindJSON(l)
	if err != nil {
		utils.Logger.Error("like action error,invalid param",
			zap.Error(err))
		str.Response(c, str.ErrInvalidParam, str.Empty, nil)
		return
	}
	l.UserId, err = utils.GetUserId(c)
	if err != nil {
		utils.Logger.Warn("user not log in,but want to like action",
			zap.Error(err))
		str.Response(c, err, str.Empty, nil)
		return
	}
	l.Url = c.Request.URL.RawPath
	_, err = client.LikeAction(c, &likePb.LikeActionRequest{
		UserId:     l.UserId,
		SourceId:   l.SourceId,
		SourceType: l.SourceType,
		ActionTye:  l.ActionType,
		Url:        l.Url,
	})
	if err != nil {
		utils.Logger.Error("like action service error",
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
	userId, err := utils.GetUserId(c)
	if err != nil {
		utils.Logger.Warn("user not log in,but want to list like post",
			zap.Error(err))
		str.Response(c, err, str.Empty, nil)
		return
	}
	resp, err := client.LikeList(c, &likePb.LikeListRequest{
		UserId: userId,
	})
	if err != nil {
		utils.Logger.Error("like list service error",
			zap.Error(err),
			zap.Int64("userId", userId))
		str.Response(c, err, str.Empty, nil)
		return
	}
	str.Response(c, nil, "posts", resp.Posts)
	return
}
