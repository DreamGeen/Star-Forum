package redis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"star/constant/settings"
)

var Client *redis.Client

func init() {
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d",
			settings.Conf.RedisHost,
			settings.Conf.RedisPort),
		Password: settings.Conf.RedisPassword, // 密码
		DB:       settings.Conf.RedisDb,       // 数据库
		PoolSize: settings.Conf.RedisPoolSize, // 连接池大小
	})
	Client = rdb
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Println("redis connect fail", err)
		panic(err)
	}
}
