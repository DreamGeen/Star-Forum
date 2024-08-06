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
