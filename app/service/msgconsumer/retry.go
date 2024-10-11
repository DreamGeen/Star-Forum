package main

import (
	"context"
	"github.com/rabbitmq/amqp091-go"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"star/app/constant/str"
	"star/app/utils/logging"
)

const (
	//最大重试次数
	maxRetries = 3
)

func sendRetryMessage(msg amqp091.Delivery, ttl int64, span trace.Span, logger *zap.Logger) {
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

	if retryCount > maxRetries {
		// 达到最大重试次数，拒绝消息且不重新入队
		if nackErr := msg.Nack(false, false); nackErr != nil {
			logger.Error("nack message error",
				zap.ByteString("message body", msg.Body),
				zap.Error(nackErr))
		}
		logger.Warn("message discarded after max retries",
			zap.Int32("retry count", retryCount))
		return
	}
	if err := msg.Ack(false); err != nil {
		logger.Error("ack msg error",
			zap.Error(err))
	}

	retryCount = retryCount + 1
	header := make(map[string]interface{}, 2)
	header["x-retry-count"] = retryCount
	header["x-delay"] = ttl
	err := channel.PublishWithContext(
		context.Background(),
		str.RetryExchange,
		msg.RoutingKey,
		false,
		false,
		amqp091.Publishing{
			Headers:      header,
			Body:         msg.Body,
			Type:         "text/plain",
			DeliveryMode: amqp091.Persistent,
		},
	)
	if err != nil {
		logger.Error("send retry message fail",
			zap.Error(err),
			zap.Any("msg", msg))
		logging.SetSpanError(span, err)
	}

}
