package RabbitMQ

import (
	"fmt"
	"github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
	"star/utils"
)

func PublishStarEvent(commentId int64) error {
	if rabbitMQConn == nil {
		err := fmt.Errorf("RabbitMQ 连接未初始化")
		utils.Logger.Error("MQ点赞生产通道获取失败", zap.Error(err))
		return err
	}

	ch, err := rabbitMQConn.Channel()
	if err != nil {
		utils.Logger.Error("MQ点赞生产通道获取失败", zap.Error(err))
		return err
	}

	q, err := ch.QueueDeclare(
		"comment_star", // 队列名称
		true,           // 是否持久化
		false,          // 是否在消费者断开连接时自动删除队列
		false,          // 是否独占队列
		false,          // 是否非阻塞模式
		nil,            // 其他参数
	)
	if err != nil {
		utils.Logger.Error("MQ声明点赞队列失败", zap.Error(err))
		return err
	}

	body := fmt.Sprintf("star:%d", commentId)
	err = ch.Publish(
		"",     // 交换机
		q.Name, // 路由键，这里使用队列名称
		false,  // 是否强制模式
		false,  // 是否立即模式
		amqp091.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	if err != nil {
		utils.Logger.Error("MQ发布点赞消息失败", zap.Error(err))
		return err
	}

	utils.Logger.Info("成功发布点赞消息", zap.String("message", body))
	return nil
}
