package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"net/http"
	"os"
	"star/app/constant/str"
	"star/app/utils/logging"
	"star/app/utils/request"
)

const (
	maxMessageSize = 1024 * 2
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func Run(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logging.Logger.Error("upgrader connection failed", zap.Error(err))
		str.Response(c, err, str.Empty, nil)
		return
	}
	//设置最大读取消息大小
	conn.SetReadLimit(maxMessageSize)
	userId, err := request.GetUserId(c)
	if err != nil {
		logging.Logger.Error("get user id failed", zap.Error(err))
		str.Response(c, err, str.Empty, nil)
		return
	}
	instanceId, err := getHostnameInstanceID()
	if err != nil {
		logging.Logger.Error("get hostname instance id failed", zap.Error(err))
		str.Response(c, str.ErrServiceBusy, str.Empty, nil)
		return
	}
	if err := SaveServiceId(userId, instanceId); err != nil {
		logging.Logger.Error("save service id failed", zap.Error(err))
		str.Response(c, err, str.Empty, nil)
		return
	}
	client := NewClient(conn, userId)

	go client.read()
	go client.write()
	manager.Connect <- client
}

// 获取主机名作为实例 ID
func getHostnameInstanceID() (string, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return "", err
	}
	return hostname, nil
}
