package redis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"star/constant/str"
	"star/models"
	"strconv"
	"time"
)

func SaveMessage(message *models.PrivateMessage) (err error) {
	key := getMessageKey(message)
	_, err = Client.ZAdd(context.Background(), key,
		redis.Z{
			Score:  float64(message.SendTime.Unix()),
			Member: message},
	).Result()
	return err
}

func BitchSaveMessage(messages []*models.PrivateMessage) error {
	pipe := Client.TxPipeline()
	for _, message := range messages {
		key := getMessageKey(message)
		pipe.ZAdd(context.Background(), key, redis.Z{
			Score:  float64(message.SendTime.Unix()),
			Member: message,
		})
	}
	_, err := pipe.Exec(context.Background())
	if err != nil {
		zap.L().Error("bitch save message failed", zap.Error(err), zap.Any("messages", messages), zap.Int64("privateChatId", messages[0].PrivateChatId))
		return str.ErrMessageError
	}
	return nil
}

func LoadMessage(senderId int64, recipientId int64, lastMsgTime time.Time, limit int64) ([]*models.PrivateMessage, error) {
	key := getMessageKey(&models.PrivateMessage{
		SenderId:    senderId,
		RecipientId: recipientId,
	})
	messages := make([]*models.PrivateMessage, 0, limit-1)
	err := Client.ZRevRangeByScore(context.Background(), key,
		&redis.ZRangeBy{
			Min:    "-inf",
			Max:    strconv.Itoa(int(lastMsgTime.Unix())),
			Offset: 0,
			Count:  limit,
		}).ScanSlice(&messages)
	if err != nil {
		zap.L().Error("redis ZRevRangeByScore err:", zap.Error(err))
		return nil, err
	}
	return messages, nil
}

func RemoveMessage(senderId int64, recipientId int64, start int64, limit int64) error {
	key := getMessageKey(&models.PrivateMessage{
		SenderId:    senderId,
		RecipientId: recipientId,
	})
	_, err := Client.ZRemRangeByRank(context.Background(), key, start, start+limit-1).Result()
	if err != nil {
		zap.L().Error("redis ZRemRangeByRank err:", zap.Error(err))
		return err
	}
	return nil
}

func GetMessageLength(senderId int64, recipientId int64) int {
	key := getMessageKey(&models.PrivateMessage{
		SenderId:    senderId,
		RecipientId: recipientId,
	})
	length, _ := Client.ZCard(context.Background(), key).Result()
	return int(length)
}

func getMessageKey(message *models.PrivateMessage) string {
	chat := models.GetPrivateChat(message)
	return fmt.Sprintf("chat:%d_%d", chat.User1Id, chat.User2Id)
}
