package rabbitMQ

import (
	"log"
	"star/app/storage/mysql"
	"strconv"
	"strings"
)

// ConsumeStarEvents 消费点赞消息
func ConsumeStarEvents() {
	if rabbitMQConn == nil {
		log.Println("rabbitMQ 连接未初始化")
		return
	}

	ch, err := rabbitMQConn.Channel()
	if err != nil {
		log.Println("MQ点赞消费通道获取失败", err)
		return
	}

	msgs, err := ch.Consume(
		"comment_star", // 队列名
		"",             // 消费者标签
		true,           // 自动确认消息
		false,          // 是否独占队列
		false,          // 是否非本地消息
		false,          // 是否非阻塞模式
		nil,            // 其他参数
	)
	if err != nil {
		log.Println("MQ点赞消息接收失败", err)
	}

	go func() {
		for d := range msgs {
			body := string(d.Body)
			if body == "heartbeat" {
				// 心跳消息
				log.Println("comment_star接收到心跳消息", "websocket", body)
				continue
			}
			if strings.HasPrefix(body, "star:") {
				commentIdStr := body[5:]
				commentId, err := strconv.ParseInt(commentIdStr, 10, 64)
				if err != nil {
					log.Println("MQ解析消息体失败", err)
					continue
				}
				if err := mysql.UpdateStar(commentId, 1); err != nil {
					log.Println("MQ更新点赞数失败", err)
				} else {
					log.Println("成功更新点赞数", "commentId", commentId)
				}
			} else {
				log.Println("MQ消息体格式错误", "body", body)
			}
		}
	}()
}

// ConsumeDeleteCommentEvents 消费评论删除事件
func ConsumeDeleteCommentEvents() {
	if rabbitMQConn == nil {
		log.Println("rabbitMQ 连接未初始化")
		return
	}

	ch, err := rabbitMQConn.Channel()
	if err != nil {
		log.Println("MQ评论删除消费通道获取失败", err)
		return
	}

	msgs, err := ch.Consume(
		"comment_delete", // 队列名
		"",               // 消费者标签
		true,             // 自动确认消息
		false,            // 是否独占队列
		false,            // 是否非本地消息
		false,            // 是否非阻塞模式
		nil,              // 其他参数
	)
	if err != nil {
		log.Println("MQ评论删除消息接收失败", err)
		return
	}

	go func() {
		for d := range msgs {
			body := string(d.Body)
			if body == "heartbeat" {
				// 心跳消息
				log.Println("comment_delete接收到心跳消息", "websocket", body)
				continue
			}
			if strings.HasPrefix(body, "delete-comment:") {
				commentIdStr := body[15:]
				commentId, err := strconv.ParseInt(commentIdStr, 10, 64)
				if err != nil {
					log.Println("MQ解析删除评论消息体失败", err)
					continue
				}
				if err := mysql.DeleteComment(commentId); err != nil {
					log.Println("MQ处理删除评论失败", err)
				} else {
					log.Println("成功处理删除评论", "commentId", commentId)
				}
			} else {
				log.Println("MQ删除评论消息体格式错误", "body", body)
			}
		}
	}()
}

//// ConsumeCommentEvents 从RabbitMQ队列中消费评论事件
//func ConsumeCommentEvents() {
//	if rabbitMQConn == nil {
//		logger.CommentLogger.Error("rabbitMQ 连接未初始化")
//		return
//	}
//
//	ch, err := rabbitMQConn.Channel()
//	if err != nil {
//		logger.CommentLogger.Error("MQ评论消费通道获取失败", zap.Error(err))
//		return
//	}
//
//	msgs, err := ch.Consume(
//		"comment_post", // 队列名
//		"",             // 消费者标签
//		false,          // 自动确认模式
//		false,          // 是否独占队列
//		false,          // 是否非本地消息
//		false,          // 是否非阻塞模式
//		nil,            // 其他参数
//	)
//	if err != nil {
//		logger.CommentLogger.Error("MQ评论消息接收失败", zap.Error(err))
//	}
//
//	go func() {
//		for d := range msgs {
//			body := string(d.Body)
//			if body == "heartbeat" {
//				// 心跳消息
//				logger.CommentLogger.Info("comment_post接收到心跳消息", zap.String("websocket", body))
//			} else {
//				var comment models.Comment
//				if err := json.Unmarshal(d.Body, &comment); err != nil {
//					logger.CommentLogger.Error("MQ反序列化评论事件失败", zap.Error(err))
//					continue
//				}
//
//				if err := mysql.CreateComment(&comment); err != nil {
//					logger.CommentLogger.Error("MQ发布评论失败", zap.Error(err))
//				}
//
//				d.Ack(false)
//			}
//		}
//	}()
//}
