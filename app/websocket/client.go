package websocket

import (
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"star/models"
	"time"
)

type Client struct {
	userId int64
	send   chan *models.Message //储存发送信息
	conn   *websocket.Conn      //websocket连接
}

func NewClient(conn *websocket.Conn, userId int64) *Client {
	return &Client{
		userId: userId,
		send:   make(chan *models.Message, 256),
		conn:   conn,
	}
}

func (c *Client) Read() {
	for {
		_, content, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				zap.L().Error("error: %v", zap.Any("err", err))
				break
			}
		}
		message := &models.Message{
			RecipientId: c.userId,
			Content:     string(content),
			SendTime:    time.Now(),
		}
		c.send <- message
	}

}
