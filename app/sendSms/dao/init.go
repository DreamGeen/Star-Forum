package redis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"star/constant/settings"
)

var rdb *redis.Client

func Init() (err error) {
	rdb = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d",
			settings.Conf.RedisHost,
			settings.Conf.RedisPort),
		Password: settings.Conf.RedisPassword, // 密码
		DB:       settings.Conf.RedisDb,       // 数据库
		PoolSize: settings.Conf.RedisPoolSize, // 连接池大小
	})
	_, err = rdb.Ping(context.Background()).Result()
	if err != nil {
		return err
	}
	return nil
}
func Close() {
	_ = rdb.Close()
}
