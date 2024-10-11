package main

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"star/app/constant/str"
	"star/app/storage/redis"
	"star/app/utils/logging"
	"sync"
)

type ClientManager struct {
	ClientMap  map[int64]*Client
	Lock       sync.RWMutex
	Connect    chan *Client
	Disconnect chan *Client
}

func NewManager() *ClientManager {
	return &ClientManager{
		ClientMap:  make(map[int64]*Client),
		Connect:    make(chan *Client, 10000),
		Disconnect: make(chan *Client, 10000),
	}
}

func (m *ClientManager) Run() {
	for {
		select {
		case client := <-m.Connect:
			m.EventConnect(client)
		case client := <-m.Disconnect:
			m.EventDisConnect(client)
		}
	}
}

// SaveServiceId 将websocket连接的serviceId储存到redis中
func SaveServiceId(userId int64, instanceId string) error {
	//将ServiceId保存到redis中去
	key := fmt.Sprintf("websocket:%d", userId)
	serviceId := "messageService:" + instanceId
	err := redis.Client.SetNX(context.Background(), key, serviceId, 0).Err()
	if err != nil {
		logging.Logger.Error("redis.SetNX serviceId failed", zap.Error(err))
		return str.ErrServiceBusy
	}
	return nil
}

func (m *ClientManager) EventConnect(client *Client) {
	m.Lock.Lock()
	defer m.Lock.Unlock()

	m.ClientMap[client.userId] = client
}

func (m *ClientManager) EventDisConnect(client *Client) {
	m.Lock.Lock()
	defer m.Lock.Unlock()

	_ = client.conn.Close()
	close(client.send)
	delete(m.ClientMap, client.userId)
}

func (m *ClientManager) GetClient(userId int64) (*Client, bool) {
	m.Lock.RLock()
	defer m.Lock.RUnlock()

	if client, ok := m.ClientMap[userId]; ok {
		return client, true
	}
	return nil, false
}
