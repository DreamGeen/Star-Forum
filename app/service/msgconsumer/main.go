package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
	"star/app/constant/str"
	"star/app/extra/tracing"
	"star/app/models"
	"star/app/storage/mysql"
	"star/app/utils/logging"
	"star/app/utils/rabbitmq"
	"star/app/utils/snowflake"
	"sync"
)

var conn *amqp091.Connection
var channel *amqp091.Channel

func failOnError(err error, msg string) {
	if err != nil {
		zap.L().Error(msg, zap.Error(err))
	}
}
func closeMQ() {
	if err := conn.Close(); err != nil {
		zap.L().Error("close rabbitmq conn error", zap.Error(err))
		panic(err)
	}
	if err := channel.Close(); err != nil {
		zap.L().Error("close rabbitmq channel error", zap.Error(err))
		panic(err)
	}
}

func main() {
	//雪花算法初始化
	if err := snowflake.Init(1); err != nil {
		panic(err)
	}
	//连接消息队列
	var err error
	conn, err = amqp091.Dial(rabbitmq.ReturnRabbitmqUrl())
	failOnError(err, "Failed to connect to RabbitMQ")

	channel, err = conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer closeMQ()

	err = channel.ExchangeDeclare(str.MessageExchange,
		"topic",
		false, false, false, false,
		nil)
	failOnError(err, "Failed to declare an exchange")
	err = channel.ExchangeDeclare(str.RetryExchange,
		"x-delayed-message",
		false, false, false, false,
		amqp091.Table{
			"x-delayed-type": "topic",
		})
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
	err = channel.QueueBind(str.MessageLike, str.RoutMessageLike, str.MessageExchange, false, nil)
	failOnError(err, "Failed to bind like queue")
	err = channel.QueueBind(str.MessageLike, str.RoutMessageLike, str.RetryExchange, false, nil)
	failOnError(err, "Failed to bind like queue")

	// 绑定@提及消息
	err = channel.QueueBind(str.MessageMention, str.RoutMention, str.MessageExchange, false, nil)
	failOnError(err, "Failed to bind mention queue")
	err = channel.QueueBind(str.MessageMention, str.RoutMention, str.RetryExchange, false, nil)
	failOnError(err, "Failed to bind mention queue")

	// 绑定回复消息
	err = channel.QueueBind(str.MessageReply, str.RoutMention, str.MessageExchange, false, nil)
	failOnError(err, "Failed to bind reply queue")
	err = channel.QueueBind(str.MessageReply, str.RoutMention, str.RetryExchange, false, nil)
	failOnError(err, "Failed to bind reply queue")

	// 绑定系统通知
	err = channel.QueueBind(str.MessageSystem, str.RoutSystem, str.MessageExchange, false, nil)
	failOnError(err, "Failed to bind system queue")
	err = channel.QueueBind(str.MessageSystem, str.RoutSystem, str.RetryExchange, false, nil)
	failOnError(err, "Failed to bind system queue")

	// 绑定私信消息
	err = channel.QueueBind(str.MessagePrivateMsg, str.RoutPrivateMsg, str.MessageExchange, false, nil)
	failOnError(err, "Failed to bind private message queue")
	err = channel.QueueBind(str.MessagePrivateMsg, str.RoutPrivateMsg, str.RetryExchange, false, nil)
	failOnError(err, "Failed to bind private message queue")
	var wg sync.WaitGroup
	wg.Add(1)

	go savePrivateMessage()

	go saveSystemMessage()

	go saveMentionMessage()

	go saveReplyMessage()

	go saveLikeMessage()
	wg.Wait()
}

func savePrivateMessage() {
	delivery, err := channel.Consume(str.MessagePrivateMsg,
		str.Empty, false, false, false, false, nil)
	failOnError(err, "Failed to register a privateMsg consumer")
	handleMessage(delivery, func(message interface{}) error {
		return mysql.InsertPrivateMsg(message.(*models.PrivateMessage))
	}, str.MessagePrivateMsg)
}

func saveSystemMessage() {
	delivery, err := channel.Consume(str.MessageSystem,
		str.Empty, false, false, false, false, nil)
	failOnError(err, "Failed to register a privateMsg consumer")
	handleMessage(delivery, func(message interface{}) error {
		return mysql.InsertSystemMsg(message.(*models.SystemMessage))
	}, str.MessageSystem)
}
func saveLikeMessage() {
	delivery, err := channel.Consume(str.MessageLike,
		str.Empty, false, false, false, false, nil)
	failOnError(err, "Failed to register a like consumer")
	handleMessage(delivery, func(message interface{}) error {
		return mysql.UpdateLikeMessage(message.(*models.RemindMessage))
	}, str.MessageLike)
}
func saveMentionMessage() {
	delivery, err := channel.Consume(str.MessageMention,
		str.Empty, false, false, false, false, nil)
	failOnError(err, "Failed to register a mention consumer")
	handleMessage(delivery, func(message interface{}) error {
		return mysql.InsertMentionMessage(message.(*models.RemindMessage))
	}, str.MessageMention)
}

func saveReplyMessage() {
	delivery, err := channel.Consume(str.MessageReply,
		str.Empty, false, false, false, false, nil)
	failOnError(err, "Failed to register a reply consumer")
	handleMessage(delivery, func(message interface{}) error {
		return mysql.InsertReplyMessage(message.(*models.RemindMessage))
	}, str.MessageReply)
}

func getFuncNewInstance(msgType string) func() interface{} {
	switch msgType {
	case str.MessagePrivateMsg:
		return func() interface{} {
			return &models.PrivateMessage{}
		}
	case str.MessageSystem:
		return func() interface{} {
			return &models.SystemMessage{}
		}
	case str.MessageMention, str.MessageLike, str.MessageReply:
		return func() interface{} {
			return &models.RemindMessage{}
		}
	}
	return nil
}

func handleMessage(delivery <-chan amqp091.Delivery, updateFunc func(interface{}) error, msgType string) {
	getInstance := getFuncNewInstance(msgType)
	for msg := range delivery {
		ctx := rabbitmq.ExtractAMQPHeaders(context.Background(), msg.Headers)

		ctx, span := tracing.Tracer.Start(ctx, "MessageSendService")
		logger := logging.LogServiceWithTrace(span, "MessageSend")
		message := getInstance()

		// 反序列化消息体
		if err := json.Unmarshal(msg.Body, message); err != nil {
			logger.Error(fmt.Sprintf("unmarshal %s message error", msgType),
				zap.ByteString("message body", msg.Body),
				zap.Error(err))
			//序列化失败，拒绝消息且不重新入队
			if nackErr := msg.Nack(false, false); nackErr != nil {
				logging.Logger.Error("nack message error",
					zap.ByteString("message body", msg.Body),
					zap.Error(nackErr))
			}
			logging.SetSpanError(span, err)
			span.End()
			continue
		}

		// 更新消息到数据库
		if err := updateFunc(message); err != nil {
			logger.Error(fmt.Sprintf("insert %s message error", msgType),
				zap.Error(err),
				zap.Any("message", message))
			sendRetryMessage(msg, 100, span, logger)
			logging.SetSpanError(span, err)
			span.End()
			continue
		}

		// 成功处理消息，确认 (ack)
		if err := msg.Ack(false); err != nil {
			logger.Error(fmt.Sprintf("ack %s message error", msgType),
				zap.Error(err))
			logging.SetSpanError(span, err)
		}
		span.End()
	}
}
