package httpHandler

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"star/app/gateway/client"
	"star/constant/str"
	"star/proto/message/messagePb"
	"star/utils"
)

func ListMessageCountHandler(c *gin.Context) {
	userId, err := utils.GetUserId(c)
	if err != nil {
		utils.Logger.Error("get user id error", zap.Error(err))
		str.Response(c, err, str.Empty, nil)
		return
	}
	resp, err := client.ListMessageCount(c,
		&messagePb.ListMessageCountRequest{
			UserId: userId,
		})
	if err != nil {
		str.Response(c, err, str.Empty, nil)
		return
	}
	str.Response(c, nil, "count", resp.Count)
}

func SendSystemHandler(c *gin.Context) {
	req := &messagePb.SendSystemMessageRequest{
		ManagerId:   1820019310731464704,
		RecipientId: 0,
		Type:        "all",
		Title:       "Test",
		Content:     "test send system message",
	}
	resp, err := client.SendSystemMessage(c, req)
	if err != nil {
		utils.Logger.Error("send system message error", zap.Error(err))
		str.Response(c, err, str.Empty, nil)
		return
	}
	str.Response(c, nil, "success", resp)

}
