package mq

import "github.com/rabbitmq/amqp091-go"

func SendMessage(queueName string, body []byte) (err error) {
	channel, _ := mq.Channel()

	q, _ := channel.QueueDeclare(queueName, true, false, false, false, nil)
	return channel.Publish("", q.Name, false, false, amqp091.Publishing{
		DeliveryMode: amqp091.Persistent,
		ContentType:  "text/plain",
		Body:         body,
	})
}

// ConsumeMessage  从mq到mysql
func ConsumeMessage(queueName string) (msg <-chan amqp091.Delivery, err error) {
	channel, _ := mq.Channel()

	q, err := channel.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		return nil, err
	}
	err = channel.Qos(1, 0, false)
	return channel.Consume(q.Name, "", false, false, false, false, nil)
}
