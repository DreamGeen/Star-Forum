package httpHandler

import (
	"github.com/gin-gonic/gin"
	"log"
	"star/app/gateway/client"
	"star/constant/str"
	"star/proto/message/messagePb"
	"star/utils"
)

func ListMessageCountHandler(c *gin.Context) {
	userId, err := utils.GetUserId(c)
	if err != nil {
		log.Println("get user id error", err)
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
