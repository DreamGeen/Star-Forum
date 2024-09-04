package rabbitMQ

import (
	"fmt"
	"github.com/rabbitmq/amqp091-go"
	"log"
	"star/constant/settings"
)

var rabbitMQConn *amqp091.Connection

func ConnectToRabbitMQ() error {
	// 格式化连接字符串

	connStr := fmt.Sprintf("amqp://%s:%s@%s:%d/", settings.Conf.RabbitMQConfig.Username, settings.Conf.RabbitMQConfig.Password, settings.Conf.RabbitMQConfig.Host, settings.Conf.RabbitMQConfig.Port)
	log.Println("尝试连接到 rabbitMQ", "connStr", connStr)

	// 尝试连接
	var err error
	rabbitMQConn, err = amqp091.Dial(connStr)
	if err != nil {
		log.Println("连接到 rabbitMQ 失败", err)
		return err
	}

	if err := DeclareQueues(); err != nil {
		return err
	}

	log.Println("成功连接到 rabbitMQ")
	return nil
}

func Close() {
	if rabbitMQConn != nil {
		err := rabbitMQConn.Close()
		if err != nil {
			log.Println("关闭 rabbitMQ 连接失败", err)
		} else {
			log.Println("成功关闭 rabbitMQ 连接")
		}
	}
}

// DeclareQueues 在服务启动时声明所有需要的队列
func DeclareQueues() error {
	if rabbitMQConn == nil {
		err := fmt.Errorf("rabbitMQ 连接未初始化")
		log.Println("队列声明失败", err)
		return err
	}

	ch, err := rabbitMQConn.Channel()
	if err != nil {
		log.Println("创建通道失败", err)
		return err
	}
	defer func() {
		_ = ch.Close()
	}()

	_, err = ch.QueueDeclare(
		"comment_star", // 队列名称
		true,           // 是否持久化
		false,          // 是否在消费者断开连接时自动删除队列
		false,          // 是否独占队列
		false,          // 是否非阻塞模式
		nil,            // 其他参数
	)
	if err != nil {
		return err
	}

	_, err = ch.QueueDeclare(
		"comment_delete", // 队列名称
		true,             // 是否持久化
		false,            // 是否在消费者断开连接时自动删除队列
		false,            // 是否独占队列
		false,            // 是否非阻塞模式
		nil,              // 其他参数
	)
	if err != nil {
		return err
	}

	//_, err = ch.QueueDeclare(
	//	"comment_post", // 队列名称
	//	true,           // 是否持久化
	//	false,          // 是否在消费者断开连接时自动删除队列
	//	false,          // 是否独占队列
	//	false,          // 是否非阻塞模式
	//	nil,            // 其他参数
	//)
	//if err != nil {
	//	return err
	//}

	return nil
}
