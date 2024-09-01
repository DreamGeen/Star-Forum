package models

import (
	"github.com/gorilla/websocket"
	"star/models"
)

type Hub struct {
	clients    map[*Client]bool          //已注册的客户端
	register   chan *Client              //注册客户端的请求
	unregister chan *Client              // 删除客户端的请求
	broadcast  chan *models.GroupMessage //广播
}

type Client struct {
	hub *Hub

	user *models.User

	send chan *models.GroupMessage //储存发送信息

	conn *websocket.Conn //websocket连接
}
