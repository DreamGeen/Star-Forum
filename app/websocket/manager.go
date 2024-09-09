package websocket

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"os"
	"star/app/storage/redis"
	"star/constant/str"
)

// 获取主机名作为实例 ID
func getHostnameInstanceID() (string, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return "", err
	}
	return hostname, nil
}

// SaveServiceId 将websocket连接的serviceId储存到redis中
func SaveServiceId(userId int64, instanceId string) error {
	//将ServiceId保存到redis中去
	key := fmt.Sprintf("websocket:%d", userId)
	serviceId := "messageService:" + instanceId
	err := redis.Client.SetNX(context.Background(), key, serviceId, 0).Err()
	if err != nil {
		zap.L().Error("redis.SetNX serviceId failed", zap.Error(err))
		return str.ErrServiceBusy
	}
	return nil
}

func SendMessage() {

}
