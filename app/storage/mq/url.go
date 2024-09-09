package mq

import (
	"fmt"
	"star/constant/settings"
)

// ReturnRabbitmqUrl 返回连接rabbitmq的url
func ReturnRabbitmqUrl() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%d/",
		settings.Conf.Username,
		settings.Conf.Password,
		settings.Conf.Host,
		settings.Conf.Port)
}
