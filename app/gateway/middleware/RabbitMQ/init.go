package RabbitMQ

import (
	"fmt"
	"github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
	"star/settings"
	"star/utils"
)

var rabbitMQConn *amqp091.Connection

func ConnectToRabbitMQ() error {
	// 格式化连接字符串
	connStr := fmt.Sprintf("amqp://%s:%s@%s:%d/", settings.Conf.RabbitMQConfig.Username, settings.Conf.RabbitMQConfig.Password, settings.Conf.RabbitMQConfig.Host, settings.Conf.RabbitMQConfig.Port)
	utils.Logger.Info("尝试连接到 RabbitMQ", zap.String("connStr", connStr))

	// 尝试连接
	var err error
	rabbitMQConn, err = amqp091.Dial(connStr)
	if err != nil {
		utils.Logger.Error("连接到 RabbitMQ 失败", zap.Error(err))
		return err
	}

	if err := DeclareQueues(); err != nil {
		return err
	}

	utils.Logger.Info("成功连接到 RabbitMQ")
	return nil
}

func Close() {
	if rabbitMQConn != nil {
		err := rabbitMQConn.Close()
		if err != nil {
			utils.Logger.Error("关闭 RabbitMQ 连接失败", zap.Error(err))
		} else {
			utils.Logger.Info("成功关闭 RabbitMQ 连接")
		}
	}
}

// DeclareQueues 在服务启动时声明所有需要的队列
func DeclareQueues() error {
	if rabbitMQConn == nil {
		err := fmt.Errorf("RabbitMQ 连接未初始化")
		utils.Logger.Error("队列声明失败", zap.Error(err))
		return err
	}

	ch, err := rabbitMQConn.Channel()
	if err != nil {
		utils.Logger.Error("创建通道失败", zap.Error(err))
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
