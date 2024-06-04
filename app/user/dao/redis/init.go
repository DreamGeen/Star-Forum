package redis

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"

	"star/settings"
)

var rdb *redis.Client

func Init() (err error) {
	rdb = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d",
			settings.Conf.RedisConfig.Host,
			settings.Conf.RedisConfig.Port),
		Password: settings.Conf.RedisConfig.Password, // 密码
		DB:       settings.Conf.RedisConfig.Db,       // 数据库
		PoolSize: settings.Conf.PoolSize,             // 连接池大小
	})
	_, err = rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Println("连接redis失败")
		return err
	}
	return nil
}
func Close() {
	_ = rdb.Close()
}
