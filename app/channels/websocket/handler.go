package websocket

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"math"
	"star/app/storage/mysql"
	"star/app/storage/redis"
	"star/models"
	"star/utils"

	"log"
	"net/http"
	"time"
)

const (
	// 允许向对方写入消息的时间。
	writeWait = 10 * time.Second

	// 允许从对方读取下一条 pong 消息的时间。
	pongWait = 60 * time.Second

	// 向与此时间段对方发送 ping。必须小于 pongWait。
	pingPeriod = (pongWait * 9) / 10

	//允许的最大消息大小。
	maxMessageSize = 512

	//将redis存储的消息储存到mysql的时间
	savePeriod = 30 * time.Second
	//redis消息的最大储存数量
	maxSaveSize = 2
	// 初始批量大小
	initialBatchSize = 2
	// 最大批量大小
	maxBatchSize = 500
)

type Client struct {
	hub *Hub

	user *models.User

	send chan *models.GroupMessage //储存发送信息

	conn *websocket.Conn //websocket连接
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (c *Client) read(communityId int64, communityIdStr string) {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	//设置读取大小限制
	c.conn.SetReadLimit(maxMessageSize)
	//设置初始读取截止时间
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	//设置 WebSocket 连接的 pong 消息处理函数
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
			}
			break
		}
		//消息处理
		chatMsg := &models.GroupMessage{
			Content:     string(msg),
			UserName:    c.user.Username,
			Img:         c.user.Img,
			SendTime:    time.Now(),
			ChatId:      utils.GetID(),
			SdUserId:    c.user.UserId,
			CommunityId: communityId,
		}
		//将消息储存到redis中
		if err := redis.SaveMsg(communityIdStr, chatMsg); err != nil {
			log.Println("redis.SaveMsg err:", err)
		}
		c.hub.broadcast <- chatMsg
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
				//hub已经关闭了send通道
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			messageJson, _ := json.Marshal(message)
			//创建编写器
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				fmt.Println("NextWriter err:", err)
			}
			w.Write(messageJson)
			//将队列的聊天消息添加到当前 websocket 消息中。
			n := len(c.send)
			for i := 0; i < n; i++ {
				msgJson, _ := json.Marshal(<-c.send)
				w.Write(msgJson)
			}
			//关闭编写器
			if err := w.Close(); err != nil {
				fmt.Println("Close err:", err)
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}

		}
	}

}

func ServeWs(communityId int64, communityIdStr string, user *models.User, hubManager *HubManager, w http.ResponseWriter, r *http.Request) {
	//将http连接升级为websocket连接
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	//根据communityId查找hub
	hub, err := hubManager.GetHub(communityId)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{
		user: user,
		hub:  hub,
		conn: conn,
		send: make(chan *models.GroupMessage, 256)}
	client.hub.register <- client
	//加载聊天记录
	recentMsg, err := redis.LoadMsg(communityIdStr, math.MaxInt64, 50)
	if err != nil {
		log.Println("load recent msg error", err)
		return
	} else {
		for _, msg := range recentMsg {
			client.send <- msg
		}
	}
	go client.write()
	go client.read(communityId, communityIdStr)
	go saveToMysql(communityIdStr)
}

func saveToMysql(communityIdStr string) {
	ticker := time.NewTicker(savePeriod)
	saveQueue := make(chan []*models.GroupMessage, 10)
	defer func() {
		ticker.Stop()
		close(saveQueue)
	}()
	go asynSavetoMysql(saveQueue)
	for {
		select {
		case <-ticker.C:
			length := redis.GetLength(communityIdStr)
			if length < maxSaveSize {
				break
			}
			batchPrepareSize := length - maxSaveSize
			batchSize := min(initialBatchSize, batchPrepareSize)
			start := int64(0)
			for batchSize > 0 {
				messages, err := redis.GetMsg(communityIdStr, start, int64(batchSize))
				if err != nil {
					log.Println("redis.GetMsg err:", err)
					continue
				}
				err = redis.RemoveMsg(communityIdStr, start, int64(batchSize))
				if err != nil {
					log.Println("redis.RemoveMsg err:", err)
				}
				if len(messages) > 0 {
					saveQueue <- messages
				}
				batchPrepareSize = batchPrepareSize - batchSize
				batchSize = min(batchPrepareSize, maxBatchSize)

			}

		default:
			time.Sleep(10 * time.Millisecond)
		}
	}

}

func asynSavetoMysql(saveQueue chan []*models.GroupMessage) {
	for {
		select {
		case messages := <-saveQueue:
			if err := mysql.SaveBathMsg(messages); err != nil {
				log.Println("mysql.SaveBathMsg err:", err)
			}
		}
	}
}
