package mq

import (
	"fmt"
	"github.com/rabbitmq/amqp091-go"
	"log"
	"star/constant/settings"
)

var mq *amqp091.Connection

// Init 初始化消息队列
func Init() error {
	//连接rabbitmq
	amqpUrl := fmt.Sprintf("amqp://%s:%s@%s:%d/",
		settings.Conf.Username,
		settings.Conf.Password,
		settings.Conf.Host,
		settings.Conf.Port)
	conn, err := amqp091.Dial(amqpUrl)
	mq = conn
	if err != nil {
		log.Println("开启rabbitmq连接失败", err)
		return err
	}
	return nil
}

// CloseMq  关闭消息队列
func CloseMq() {
	_ = mq.Close()
}
