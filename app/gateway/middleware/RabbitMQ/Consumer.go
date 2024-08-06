package RabbitMQ

import (
	"encoding/json"
	"go.uber.org/zap"
	"star/app/comment/dao/mysql"
	"star/models"
	"star/utils"
	"strconv"
	"strings"
)

// ConsumeStarEvents 消费点赞消息
func ConsumeStarEvents() {
	if rabbitMQConn == nil {
		utils.Logger.Error("RabbitMQ 连接未初始化")
		return
	}

	ch, err := rabbitMQConn.Channel()
	if err != nil {
		utils.Logger.Error("MQ点赞消费通道获取失败", zap.Error(err))
		return
	}

	msgs, _ := ch.Consume(
		"comment_star", // 队列名
		"",             // 消费者标签
		true,           // 自动确认消息
		false,          // 是否独占队列
		false,          // 是否非本地消息
		false,          // 是否非阻塞模式
		nil,            // 其他参数
	)

	go func() {
		for d := range msgs {
			body := string(d.Body)
			if strings.HasPrefix(body, "star:") {
				commentIdStr := body[5:]
				commentId, err := strconv.ParseInt(commentIdStr, 10, 64)
				if err != nil {
					utils.Logger.Error("MQ解析消息体失败", zap.Error(err))
					continue
				}
				if err := mysql.UpdateStar(commentId, 1); err != nil {
					utils.Logger.Error("MQ更新点赞数失败", zap.Error(err))
				} else {
					utils.Logger.Info("成功更新点赞数", zap.Int64("commentId", commentId))
				}
			} else {
				utils.Logger.Error("MQ消息体格式错误", zap.String("body", body))
			}
		}
	}()
}

// ConsumeCommentEvents 从RabbitMQ队列中消费评论事件
func ConsumeCommentEvents() {
	if rabbitMQConn == nil {
		utils.Logger.Error("RabbitMQ 连接未初始化")
		return
	}

	ch, err := rabbitMQConn.Channel()
	if err != nil {
		utils.Logger.Error("MQ评论消费通道获取失败", zap.Error(err))
		return
	}

	msgs, _ := ch.Consume(
		"comment_post", // 队列名
		"",             // 消费者标签
		false,          // 自动确认模式
		false,          // 是否独占队列
		false,          // 是否非本地消息
		false,          // 是否非阻塞模式
		nil,            // 其他参数
	)

	go func() {
		for d := range msgs {
			var comment models.Comment
			if err := json.Unmarshal(d.Body, &comment); err != nil {
				utils.Logger.Error("MQ反序列化评论事件失败", zap.Error(err))
				continue
			}

			if err := mysql.CreateComment(&comment); err != nil {
				utils.Logger.Error("MQ发布评论失败", zap.Error(err))
			}

			d.Ack(false)
		}
	}()
}
