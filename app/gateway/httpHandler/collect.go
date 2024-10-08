package httpHandler

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"star/app/gateway/client"
	"star/app/gateway/models"
	"star/constant/str"
	"star/proto/collect/collectPb"
	"star/utils"
)

func CollectActionHandler(c *gin.Context) {
	co := new(models.CollectAction)
	err := c.ShouldBindJSON(co)
	if err != nil {
		utils.Logger.Error("collect action error,invalid param",
			zap.Error(err))
		str.Response(c, str.ErrInvalidParam, str.Empty, nil)
		return
	}
	co.ActionId, err = utils.GetUserId(c)
	if err != nil {
		utils.Logger.Warn("user not log in,but want to collect action",
			zap.Error(err))
		str.Response(c, err, str.Empty, nil)
		return
	}
	_, err = client.CollectAction(c, &collectPb.CollectActionRequest{
		ActorId:    co.ActionId,
		PostId:     co.PostId,
		ActionType: co.ActionType,
	})
	if err != nil {
		utils.Logger.Error("collect action service error",
			zap.Error(err),
			zap.Int64("userId", co.ActionId),
			zap.Int64("postId", co.PostId),
			zap.Uint32("actionType", co.ActionType))
		str.Response(c, err, str.Empty, nil)
		return
	}
	str.Response(c, nil, str.Empty, nil)
	return
}

func CollectListHandler(c *gin.Context) {
	userId, err := utils.GetUserId(c)
	if err != nil {
		utils.Logger.Warn("user not log in,but want to list collect post",
			zap.Error(err))
		str.Response(c, err, str.Empty, nil)
		return
	}
	resp, err := client.CollectList(c, &collectPb.CollectListRequest{
		ActorId: userId,
	})
	if err != nil {
		utils.Logger.Error("collect list service error",
			zap.Error(err),
			zap.Int64("userId", userId))
		str.Response(c, err, str.Empty, nil)
		return
	}
	str.Response(c, nil, "posts", resp.Posts)
	return

}
