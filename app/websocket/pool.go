package websocket

import (
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"net/http"
	"star/constant/str"
	"sync"
)

// websocket连接池
var connPool sync.Map

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// SaveClient 储存websocket Client
func SaveClient(conn *websocket.Conn, userId int64) {
	client := NewClient(conn, userId)
	connPool.Store(userId, client)
}

// GetClient 获取websocket Client
func GetClient(userId int64) (*websocket.Conn, error) {
	conn, ok := connPool.Load(userId)
	if !ok {
		zap.L().Error("用户连接不存在", zap.Int64("userId", userId))
		return nil, str.ErrServiceBusy
	}
	return conn.(*websocket.Conn), nil
}

// DeleteClient 删除websocket Client
func DeleteClient(userId int64) {
	connPool.Delete(userId)
}
