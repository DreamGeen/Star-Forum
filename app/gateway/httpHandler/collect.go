package httpHandler

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	str2 "star/app/constant/str"
	"star/app/gateway/client"
	"star/app/gateway/models"
	"star/app/utils/logging"
	utils2 "star/app/utils/request"
	"star/proto/collect/collectPb"
)

func CollectActionHandler(c *gin.Context) {
	co := new(models.CollectAction)
	err := c.ShouldBindJSON(co)
	if err != nil {
		logging.Logger.Error("collect action error,invalid param",
			zap.Error(err))
		str2.Response(c, str2.ErrInvalidParam, str2.Empty, nil)
		return
	}
	co.ActionId, err = utils2.GetUserId(c)
	if err != nil {
		logging.Logger.Warn("user not log in,but want to collect action",
			zap.Error(err))
		str2.Response(c, err, str2.Empty, nil)
		return
	}
	_, err = client.CollectAction(c, &collectPb.CollectActionRequest{
		ActorId:    co.ActionId,
		PostId:     co.PostId,
		ActionType: co.ActionType,
	})
	if err != nil {
		logging.Logger.Error("collect action service error",
			zap.Error(err),
			zap.Int64("userId", co.ActionId),
			zap.Int64("postId", co.PostId),
			zap.Uint32("actionType", co.ActionType))
		str2.Response(c, err, str2.Empty, nil)
		return
	}
	str2.Response(c, nil, str2.Empty, nil)
	return
}

func CollectListHandler(c *gin.Context) {
	userId, err := utils2.GetUserId(c)
	if err != nil {
		logging.Logger.Warn("user not log in,but want to list collect feed",
			zap.Error(err))
		str2.Response(c, err, str2.Empty, nil)
		return
	}
	resp, err := client.CollectList(c, &collectPb.CollectListRequest{
		ActorId: userId,
	})
	if err != nil {
		logging.Logger.Error("collect list service error",
			zap.Error(err),
			zap.Int64("userId", userId))
		str2.Response(c, err, str2.Empty, nil)
		return
	}
	str2.Response(c, nil, "posts", resp.Posts)
	return

}
