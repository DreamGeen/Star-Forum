package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/rabbitmq/amqp091-go"
	"go-micro.dev/v4"
	"go.uber.org/zap"
	"log"
	"star/app/storage/cached"
	"star/app/storage/mq"
	"star/app/storage/mysql"
	"star/constant/str"
	"star/models"
	"star/proto/message/messagePb"
	"star/proto/user/userPb"
	"star/utils"
	"time"
)

type MessageSrv struct {
}

var userService userPb.UserService

var conn *amqp091.Connection
var channel *amqp091.Channel

var MessageTypes = []string{"mention", "like", "reply", "system", "privateMsg"}

func failOnError(err error, msg string) {
	if err != nil {
		zap.L().Error(msg, zap.Error(err))
	}
}

func CloseMQ() {
	if err := conn.Close(); err != nil {
		zap.L().Error("close rabbitmq conn error", zap.Error(err))
		panic(err)
	}
	if err := channel.Close(); err != nil {
		zap.L().Error("close rabbitmq channel error", zap.Error(err))
		panic(err)
	}
}

func (m *MessageSrv) New() {
	//连接消息队列
	var err error
	conn, err = amqp091.Dial(mq.ReturnRabbitmqUrl())
	failOnError(err, "Failed to connect to RabbitMQ")

	channel, err = conn.Channel()
	failOnError(err, "Failed to open a channel")

	err = channel.ExchangeDeclare(str.MessageExchange,
		"topic",
		false, false, false, false,
		nil)
	failOnError(err, "Failed to declare an exchange")
	//声明队列
	_, err = channel.QueueDeclare(str.MessageLike,
		false, false, false, false,
		nil)
	failOnError(err, "Failed to declare a like queue")
	_, err = channel.QueueDeclare(str.MessageReply,
		false, false, false, false,
		nil)
	failOnError(err, "Failed to declare a reply queue")
	_, err = channel.QueueDeclare(str.MessageSystem,
		false, false, false, false,
		nil)
	failOnError(err, "Failed to declare a system queue")
	_, err = channel.QueueDeclare(str.MessagePrivateMsg,
		false, false, false, false,
		nil)
	failOnError(err, "Failed to declare a private_msg queue")
	_, err = channel.QueueDeclare(str.MessageMention,
		false, false, false, false,
		nil)
	failOnError(err, "Failed to declare a mention queue")

	//绑定队列
	// 绑定点赞消息
	err = channel.QueueBind(str.MessageLike, str.RoutLike, str.MessageExchange, false, nil)
	failOnError(err, "Failed to bind like queue")

	// 绑定@提及消息
	err = channel.QueueBind(str.MessageMention, str.RoutMention, str.MessageExchange, false, nil)
	failOnError(err, "Failed to bind mention queue")

	// 绑定回复消息
	err = channel.QueueBind(str.MessageReply, str.RoutMention, str.MessageExchange, false, nil)
	failOnError(err, "Failed to bind reply queue")

	// 绑定系统通知
	err = channel.QueueBind(str.MessageSystem, str.RoutSystem, str.MessageExchange, false, nil)
	failOnError(err, "Failed to bind system queue")

	// 绑定私信消息
	err = channel.QueueBind(str.MessagePrivateMsg, str.RoutPrivateMsg, str.MessageExchange, false, nil)
	failOnError(err, "Failed to bind private message queue")

	//创建一个用户微服务客户端
	userMicroService := micro.NewService(micro.Name(str.UserServiceClient))
	userService = userPb.NewUserService(str.UserService, userMicroService.Client())

}

func (m *MessageSrv) ListMessageCount(ctx context.Context, req *messagePb.ListMessageCountRequest, resp *messagePb.ListMessageCountResponse) error {

	//获取countsStr
	key := fmt.Sprintf("ListMessageCount:+%d", req.UserId)
	countsStr, err := cached.GetWithFunc(ctx, key, func(key string) (string, error) {
		counts, err := mysql.ListMessageCount(req.UserId)
		if err != nil {
			return "", err
		}
		countsJson, err := json.Marshal(counts)
		if err != nil {
			log.Println("json marshal counts error:", err)
			return "", str.ErrMessageError
		}
		return string(countsJson), nil

	})
	if err != nil {
		return err
	}
	//解析countsStr
	counts := new(models.Counts)
	err = json.Unmarshal([]byte(countsStr), counts)
	if err != nil {
		log.Println("json unmarshal counts error:", err)
		return str.ErrMessageError
	}
	resp.Count = &messagePb.Counts{
		PrivateMsgCount: counts.PrivateMsgCount,
		MentionCount:    counts.MentionCount,
		LikeCount:       counts.LikeCount,
		ReplyCount:      counts.ReplyCount,
		SystemCount:     counts.SystemCount,
		TotalCount:      counts.TotalCount,
	}
	return nil
}

func (m *MessageSrv) SendSystemMessage(ctx context.Context, req *messagePb.SendSystemMessageRequest, resp *messagePb.SendSystemMessageResponse) error {
	message := &models.SystemMessage{
		Id:          utils.GetID(),
		RecipientId: req.RecipientId,
		ManagerId:   req.ManagerId,
		Type:        req.Type,
		Title:       req.Title,
		Content:     req.Content,
		Status:      false,
		PublishTime: time.Now(),
	}
	body, err := json.Marshal(message)
	if err != nil {
		zap.L().Error("json marshal system message error", zap.Error(err), zap.Any("message", message))
		return str.ErrMessageError
	}
	err = channel.Publish(str.MessageExchange, str.RoutSystem,
		false, false,
		amqp091.Publishing{
			DeliveryMode: amqp091.Persistent,
			ContentType:  "text/plain",
			Body:         body,
		})
	if err != nil {
		zap.L().Error("send system message error", zap.Error(err), zap.Any("message", message))
		return str.ErrMessageError
	}
	return nil
}

func (m *MessageSrv) SendPrivateMessage(ctx context.Context, req *messagePb.SendPrivateMessageRequest, resp *messagePb.SendPrivateMessageResponse) error {
	message := &models.PrivateMessage{
		Id:          utils.GetID(),
		SenderId:    req.SenderId,
		RecipientId: req.RecipientId,
		Content:     req.Content,
		Status:      false,
		SendTime:    time.Now(),
	}
	body, err := json.Marshal(message)
	if err != nil {
		zap.L().Error("json marshal message error", zap.Error(err), zap.Any("message", message))
		return err
	}
	err = channel.Publish(str.MessageExchange, str.RoutPrivateMsg,
		false, false,
		amqp091.Publishing{
			DeliveryMode: amqp091.Persistent,
			ContentType:  "text/plain",
			Body:         body,
		},
	)
	if err != nil {
		zap.L().Error("publish message error", zap.Error(err), zap.Any("message", message))
		return err
	}
	return nil
}

func (m *MessageSrv) SendRemindMessage(ctx context.Context, req *messagePb.SendRemindMessageRequest, resp *messagePb.SendRemindMessageResponse) error {
	switch req.RemindType {
	case "like":
		if err := addRemindMessage(req, str.RoutLike); err != nil {
			zap.L().Error("add like message error", zap.Error(err), zap.Any("message", req))
			return str.ErrMessageError
		}
	case "reply":
		if err := addRemindMessage(req, str.RoutReply); err != nil {
			zap.L().Error("add reply message error", zap.Error(err), zap.Any("message", req))
			return str.ErrMessageError
		}
	case "mention":
		if err := addRemindMessage(req, str.RoutMention); err != nil {
			zap.L().Error("add mention message error", zap.Error(err), zap.Any("message", req))
			return str.ErrMessageError
		}
	}
	return nil
}

func addRemindMessage(req *messagePb.SendRemindMessageRequest, routingKey string) error {
	message := &models.RemindMessage{
		Id:          utils.GetID(),
		SourceId:    req.SourceId,
		SourceType:  req.SourceType,
		SenderId:    req.SenderId,
		RecipientId: req.RecipientId,
		Content:     req.Content,
		Url:         req.Url,
		Status:      false,
		RemindTime:  time.Now(),
	}
	body, err := json.Marshal(message)
	if err != nil {
		zap.L().Error("json marshal message error", zap.Error(err), zap.Any("message", message))
		return err
	}
	err = channel.Publish(str.MessageExchange, routingKey, false, false,
		amqp091.Publishing{
			DeliveryMode: amqp091.Persistent,
			ContentType:  "text/plain",
			Body:         body,
		})
	if err != nil {
		zap.L().Error("publish message error", zap.Error(err), zap.Any("message", message))
		return err
	}
	return nil
}

func (m *MessageSrv) SendMessage(ctx context.Context, req *messagePb.SendMessageRequest, resp *messagePb.SendMessageResponse) error {

	return nil
}
