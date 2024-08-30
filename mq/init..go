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
	amqpUrl := fmt.Sprintf("%s://%s:%s@%s:%d/",
		settings.Conf.RabbitmqAgreement,
		settings.Conf.RabbitmqUser,
		settings.Conf.RabbitmqPassword,
		settings.Conf.RabbitmqHost,
		settings.Conf.RabbitmqPort)
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
