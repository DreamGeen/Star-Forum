package rabbitMQ

import (
	"fmt"
	"github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
	logger "star/app/comment/logger"
	"star/constant/str"
	"time"
)

// PublishStarEvent 发布点赞消息
func PublishStarEvent(commentId int64) error {
	if rabbitMQConn == nil {
		err := fmt.Errorf("rabbitMQ 连接未初始化")
		logger.CommentLogger.Error("MQ点赞生产通道获取失败", zap.Error(err))
		return err
	}

	ch, err := rabbitMQConn.Channel()
	if err != nil {
		logger.CommentLogger.Error("MQ点赞生产通道获取失败", zap.Error(err))
		return err
	}
	defer func() {
		_ = ch.Close()
	}()

	body := fmt.Sprintf("star:%d", commentId)
	err = ch.Publish(
		"",             // 交换机
		"comment_star", // 路由键，使用队列名称
		false,          // 是否强制模式
		false,          // 是否立即模式
		amqp091.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	if err != nil {
		logger.CommentLogger.Error("MQ发布点赞消息失败", zap.Error(err))
		return err
	}

	logger.CommentLogger.Info("成功发布点赞消息", zap.String("message", body))
	return nil
}

// PublishDeleteEvent 发布评论删除事件到RabbitMQ
func PublishDeleteEvent(commentId int64) error {
	if rabbitMQConn == nil {
		err := fmt.Errorf("rabbitMQ 连接未初始化")
		logger.CommentLogger.Error("MQ评论删除生产通道获取失败", zap.Error(err))
		return str.ErrCommentError
	}

	ch, err := rabbitMQConn.Channel()
	if err != nil {
		logger.CommentLogger.Error("MQ评论删除生产通道获取失败", zap.Error(err))
		return str.ErrCommentError
	}
	defer func() {
		_ = ch.Close()
	}()

	// 构造消息体
	body := fmt.Sprintf("delete-comment:%d", commentId)
	err = ch.Publish(
		"",               // 交换机
		"comment_delete", // 路由键，使用队列名称
		false,            // 是否强制模式
		false,            // 是否立即模式
		amqp091.Publishing{
			ContentType:  "text/plain",
			DeliveryMode: amqp091.Persistent, // 持久化消息
			Body:         []byte(body),
		})
	if err != nil {
		logger.CommentLogger.Error("MQ发布删除评论消息失败", zap.Error(err))
		return str.ErrCommentError
	}

	logger.CommentLogger.Info("成功发布删除评论消息", zap.String("message", body))
	return nil
}

// SendHeartbeat 发送心跳消息到RabbitMQ的特定队列
func SendHeartbeat(queueName string) error {
	if rabbitMQConn == nil {
		err := fmt.Errorf("rabbitMQ 连接未初始化")
		logger.CommentLogger.Error("发送心跳消息失败", zap.Error(err))
		return err
	}

	ch, err := rabbitMQConn.Channel()
	if err != nil {
		logger.CommentLogger.Error("创建通道失败", zap.Error(err))
		return err
	}
	defer func() {
		_ = ch.Close()
	}()

	// 声明队列
	q, err := ch.QueueDeclare(
		queueName, // 队列名称
		true,      // 是否持久化
		false,     // 是否在消费者断开连接时自动删除队列
		false,     // 是否独占队列
		false,     // 是否非阻塞模式
		nil,       // 其他参数
	)
	if err != nil {
		logger.CommentLogger.Error("声明队列失败", zap.Error(err))
		return err
	}

	// 创建心跳消息体
	body := []byte("heartbeat")

	err = ch.Publish(
		"",     // 交换机
		q.Name, // 路由键
		false,  // 强制模式
		false,  // 立即模式
		amqp091.Publishing{
			ContentType: "text/plain",
			Body:        body,
		})
	if err != nil {
		logger.CommentLogger.Error("发送心跳消息失败", zap.Error(err))
		return err
	}

	logger.CommentLogger.Info("成功发送心跳消息")
	return nil
}

// StartHeartbeatTicker 开始定期发送心跳消息
func StartHeartbeatTicker(queueName string, interval time.Duration, stopCh <-chan struct{}) {
	// 每段时间都会发送一个信号
	ticker := time.NewTicker(interval)
	for {
		select {
		case <-ticker.C:
			// 发送心跳消息
			_ = SendHeartbeat(queueName)
		// 接受心跳停止信号
		case <-stopCh:
			ticker.Stop()
			return
		}
	}
}

//// PublishCommentEvent 发布评论事件到RabbitMQ
//func PublishCommentEvent(comment *models.Comment) error {
//	if rabbitMQConn == nil {
//		err := fmt.Errorf("rabbitMQ 连接未初始化")
//		logger.CommentLogger.Error("MQ评论生产通道获取失败", zap.Error(err))
//		return err
//	}
//
//	ch, err := rabbitMQConn.Channel()
//	if err != nil {
//		logger.CommentLogger.Error("MQ评论生产通道获取失败", zap.Error(err))
//		return err
//	}
//	defer func() {
//		_ = ch.Close()
//	}()
//
//	body, err := json.Marshal(comment)
//	if err != nil {
//		logger.CommentLogger.Error("序列化评论事件失败", zap.Error(err))
//		return err
//	}
//
//	err = ch.Publish(
//		"",             // 交换机
//		"comment_post", // 路由键
//		false,          // 强制模式
//		false,          // 立即模式
//		amqp091.Publishing{
//			ContentType: "application/json",
//			Body:        body,
//		})
//	if err != nil {
//		logger.CommentLogger.Error("MQ发布评论事件失败", zap.Error(err))
//		return err
//	}
//
//	logger.CommentLogger.Info("成功发布评论事件", zap.String("message", string(body)))
//	return nil
//}
