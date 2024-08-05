package redis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"star/settings"

	"go-micro.dev/v4/logger"
)

var Client *redis.Client
var Ctx = context.Background()

func Init() error {
	Client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", settings.Conf.RedisConfig.RedisHost, settings.Conf.RedisConfig.RedisPort),
		Password: settings.Conf.RedisConfig.RedisPassword,
		DB:       settings.Conf.RedisConfig.RedisDb,
	})

	_, err := Client.Ping(Ctx).Result()
	if err != nil {
		fmt.Println("连接redis失败")
		logger.Fatal(err)
	}
	return err
}

func Close() {
	_ = Client.Close()
}
