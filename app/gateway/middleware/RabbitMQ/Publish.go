package RabbitMQ

import (
	"encoding/json"
	"fmt"
	"github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
	"star/models"
	"star/utils"
)

// PublishStarEvent 发布点赞消息
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

// PublishCommentEvent 发布评论事件到RabbitMQ
func PublishCommentEvent(comment *models.Comment) error {
	if rabbitMQConn == nil {
		err := fmt.Errorf("RabbitMQ 连接未初始化")
		utils.Logger.Error("MQ评论生产通道获取失败", zap.Error(err))
		return err
	}

	ch, err := rabbitMQConn.Channel()
	if err != nil {
		utils.Logger.Error("MQ评论生产通道获取失败", zap.Error(err))
		return err
	}

	q, err := ch.QueueDeclare(
		"comment_post", // 队列名称
		true,           // 是否持久化
		false,          // 是否在消费者断开连接时自动删除队列
		false,          // 是否独占队列
		false,          // 是否非阻塞模式
		nil,            // 其他参数
	)
	if err != nil {
		utils.Logger.Error("MQ声明评论事件队列失败", zap.Error(err))
		return err
	}

	body, err := json.Marshal(comment)
	if err != nil {
		utils.Logger.Error("序列化评论事件失败", zap.Error(err))
		return err
	}

	err = ch.Publish(
		"",     // 交换机
		q.Name, // 路由键
		false,  // 强制模式
		false,  // 立即模式
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	if err != nil {
		utils.Logger.Error("MQ发布评论事件失败", zap.Error(err))
		return err
	}

	utils.Logger.Info("成功发布评论事件", zap.String("message", string(body)))
	return nil
}
