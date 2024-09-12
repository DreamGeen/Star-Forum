package test

import (
	"context"
	"encoding/json"
	"fmt"
	"star/app/storage/redis"
	"star/constant/str"
	"star/utils"
	"strconv"
	"testing"
	"time"
)

func Test_ParseTimeFormat(t *testing.T) {
	var now time.Time
	now = time.Now().UTC()
	t.Log(now)
	nowStr := now.Format(str.ParseTimeFormat)
	t.Log(nowStr)
}

type message struct {
	Id      int64  `json:"id"`
	Content string `json:"content"`
}

func Test_SaveSlice(t *testing.T) {
	if err := utils.Init(1); err != nil {
		t.Error(err)
		return
	}
	messages := make([]*message, 50)
	for i := 0; i < 50; i++ {
		msg := &message{
			Id:      utils.GetID(),
			Content: strconv.Itoa(i),
		}
		messages[i] = msg
	}
	//messageJson, err := json.Marshal(messages)
	//if err != nil {
	//	t.Error(err)
	//	return
	//}
	//err = redis.Client.Set(context.Background(), "messages", string(messageJson), 5*time.Minute).Err()
	//if err != nil {
	//	t.Error(err)
	//	return
	//}
	getMessageJson, err := redis.Client.Get(context.Background(), "messages").Result()
	if err != nil {
		t.Error(err)
		return
	}
	var msgs []*message
	err = json.Unmarshal([]byte(getMessageJson), &msgs)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(msgs)
	for _, msg := range msgs {
		fmt.Println(msg.Id, msg.Content)
	}
}
