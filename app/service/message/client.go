package main

import (
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"star/app/constant/str"
	"star/app/models"
	"star/app/service/message/service"
	"star/app/utils/logging"
	"time"
)

const (
	// 允许向对方写入消息的时间。
	writeWait = 10 * time.Second

	// 允许从对方读取下一条 pong 消息的时间。
	pongWait = 60 * time.Second

	// 向与此时间段对方发送 ping。必须小于 pongWait。
	pingPeriod = (pongWait * 9) / 10
)

type Client struct {
	userId int64
	send   chan *models.PrivateMessage //储存发送信息
	conn   *websocket.Conn             //websocket连接
}

func NewClient(conn *websocket.Conn, userId int64) *Client {
	return &Client{
		userId: userId,
		send:   make(chan *models.PrivateMessage, 256),
		conn:   conn,
	}
}

func (c *Client) read() {
	for {
		messageType, _, err := c.conn.ReadMessage()
		if err != nil {
			if messageType == -1 && websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure, websocket.CloseNoStatusReceived) {
				service.manager.Disconnect <- c
				return
			} else if messageType != websocket.PingMessage {
				return
			}
		}
	}
}

func (c *Client) write() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				//send通道已关闭
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.conn.WriteJSON(message); err != nil {
				logging.Logger.Error("write message error",
					zap.Error(err),
					zap.Int64("privateMessageId", message.Id),
					zap.Int64("privateChatId", message.PrivateChatId),
					zap.Int64("senderId", message.SenderId),
					zap.Int64("recipientId", message.RecipientId),
					zap.String("content", message.Content),
					zap.String("sendTime", message.SendTime.Format(str.ParseTimeFormat)),
					zap.Bool("status", message.Status))
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}

	}
}
