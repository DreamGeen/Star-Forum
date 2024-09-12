package main

import (
	"github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
	"star/constant/str"
	"star/utils"
)

func sendRetryMessage(msg amqp091.Delivery, ttl int64) {
	if ttl <= 0 {
		return
	}
	// 获取重试次数
	retryCount := int32(0)
	if msg.Headers == nil {
		msg.Headers = make(map[string]interface{})
	}
	if count, ok := msg.Headers["x-retry-count"].(int32); ok {
		retryCount = count
	}
	retryCount = retryCount + 1
	header := make(map[string]interface{}, 2)
	header["x-retry-count"] = retryCount
	header["x-delay"] = ttl
	err := channel.Publish(
		str.RetryExchange,
		msg.RoutingKey,
		false, false,
		amqp091.Publishing{
			Headers:      header,
			Body:         msg.Body,
			Type:         "text/plain",
			DeliveryMode: amqp091.Persistent,
		},
	)
	if err != nil {
		utils.Logger.Error("send retry message fail", zap.Error(err), zap.Any("msg", msg))
	}

}
