package redis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"star/app/models"
	"strconv"
	"time"
)

func SaveMessage(ctx context.Context, message *models.PrivateMessage) (err error) {
	key := getMessageKey(message)
	_, err = Client.ZAdd(ctx, key,
		redis.Z{
			Score:  float64(message.SendTime.Unix()),
			Member: message},
	).Result()
	return err
}

func BitchSaveMessage(ctx context.Context, messages []*models.PrivateMessage) error {
	pipe := Client.TxPipeline()
	for _, message := range messages {
		key := getMessageKey(message)
		pipe.ZAdd(ctx, key, redis.Z{
			Score:  float64(message.SendTime.Unix()),
			Member: message,
		})
	}
	_, err := pipe.Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func LoadMessage(ctx context.Context, senderId int64, recipientId int64, lastMsgTime time.Time, limit int64) ([]*models.PrivateMessage, error) {
	key := getMessageKey(&models.PrivateMessage{
		SenderId:    senderId,
		RecipientId: recipientId,
	})
	messages := make([]*models.PrivateMessage, 0, limit-1)
	err := Client.ZRevRangeByScore(ctx, key,
		&redis.ZRangeBy{
			Min:    "-inf",
			Max:    strconv.Itoa(int(lastMsgTime.Unix())),
			Offset: 0,
			Count:  limit,
		}).ScanSlice(&messages)
	if err != nil {
		return nil, err
	}
	return messages, nil
}

func RemoveMessage(ctx context.Context, senderId int64, recipientId int64, start int64, limit int64) error {
	key := getMessageKey(&models.PrivateMessage{
		SenderId:    senderId,
		RecipientId: recipientId,
	})
	_, err := Client.ZRemRangeByRank(ctx, key, start, start+limit-1).Result()
	if err != nil {
		zap.L().Error("redis ZRemRangeByRank err:", zap.Error(err))
		return err
	}
	return nil
}

func GetMessageLength(ctx context.Context, senderId int64, recipientId int64) int {
	key := getMessageKey(&models.PrivateMessage{
		SenderId:    senderId,
		RecipientId: recipientId,
	})
	length, _ := Client.ZCard(ctx, key).Result()
	return int(length)
}

func getMessageKey(message *models.PrivateMessage) string {
	chat := models.GetPrivateChat(message)
	return fmt.Sprintf("chat:%d_%d", chat.User1Id, chat.User2Id)
}
