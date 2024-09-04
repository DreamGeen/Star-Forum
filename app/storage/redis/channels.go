package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log"
	"star/models"
	"strconv"
)

func SaveMsg(communityIdStr string, message *models.GroupMessage) (err error) {
	key := "groupChat" + communityIdStr
	_, err = Client.ZAdd(context.Background(), key,
		redis.Z{
			Score:  float64(message.ChatId),
			Member: message},
	).Result()
	return err
}

func GetMsg(community string, start int64, limit int64) ([]*models.GroupMessage, error) {
	key := "groupChat" + community
	messages := make([]*models.GroupMessage, 0, limit)
	err := Client.ZRange(context.Background(), key, start, start+limit-1).ScanSlice(&messages)
	if err != nil {
		log.Println("redis zrange err:", err)
		return nil, err
	}
	return messages, nil
}

func LoadMsg(community string, lastChatIndex int64, limit int64) ([]*models.GroupMessage, error) {
	key := "groupChat" + community
	messages := make([]*models.GroupMessage, 0, limit-1)
	err := Client.ZRevRangeByScore(context.Background(), key,
		&redis.ZRangeBy{
			Min:    "-inf",
			Max:    strconv.FormatInt(lastChatIndex, 10),
			Offset: 0,
			Count:  limit,
		}).ScanSlice(&messages)
	if err != nil {
		log.Println("redis ZRevRangeByScore err:", err)
		return nil, err
	}
	return messages, nil
}

func GetLength(community string) int {
	key := "group_messages:" + community
	length, _ := Client.ZCard(context.Background(), key).Result()
	return int(length)
}

func RemoveMsg(community string, start int64, limit int64) error {
	key := "groupChat" + community
	_, err := Client.ZRemRangeByRank(context.Background(), key, start, start+limit-1).Result()
	if err != nil {
		log.Println("redis ZRemRangeByRank err:", err)
		return err
	}
	return nil
}
