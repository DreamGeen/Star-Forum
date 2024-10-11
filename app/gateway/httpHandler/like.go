package httpHandler

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	str2 "star/app/constant/str"
	"star/app/gateway/client"
	"star/app/gateway/models"
	"star/app/utils/logging"
	utils2 "star/app/utils/request"
	"star/proto/like/likePb"
)

func LikeActionHandler(c *gin.Context) {
	l := new(models.LikeAction)
	err := c.ShouldBindJSON(l)
	if err != nil {
		logging.Logger.Error("like action error,invalid param",
			zap.Error(err))
		str2.Response(c, str2.ErrInvalidParam, str2.Empty, nil)
		return
	}
	l.UserId, err = utils2.GetUserId(c)
	if err != nil {
		logging.Logger.Warn("user not log in,but want to like action",
			zap.Error(err))
		str2.Response(c, err, str2.Empty, nil)
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
		logging.Logger.Error("like action service error",
			zap.Error(err),
			zap.Int64("userId", l.UserId),
			zap.Int64("sourceId", l.SourceId),
			zap.Uint32("sourceType", l.SourceType),
			zap.String("url", l.Url))
		str2.Response(c, err, str2.Empty, nil)
		return
	}
	str2.Response(c, nil, str2.Empty, nil)
	return
}

func LikeListHandler(c *gin.Context) {
	userId, err := utils2.GetUserId(c)
	if err != nil {
		logging.Logger.Warn("user not log in,but want to list like post",
			zap.Error(err))
		str2.Response(c, err, str2.Empty, nil)
		return
	}
	resp, err := client.LikeList(c, &likePb.LikeListRequest{
		UserId: userId,
	})
	if err != nil {
		logging.Logger.Error("like list service error",
			zap.Error(err),
			zap.Int64("userId", userId))
		str2.Response(c, err, str2.Empty, nil)
		return
	}
	str2.Response(c, nil, "posts", resp.Posts)
	return
}
