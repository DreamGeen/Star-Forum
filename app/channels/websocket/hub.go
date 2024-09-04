package websocket

import "star/models"

type Hub struct {
	clients    map[*Client]bool          //已注册的客户端
	register   chan *Client              //注册客户端的请求
	unregister chan *Client              // 删除客户端的请求
	broadcast  chan *models.GroupMessage //广播
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan *models.GroupMessage),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
			}
		case messages := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- messages:
					//fmt.Println(messages)
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
