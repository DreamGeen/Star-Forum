package main

import (
	"encoding/json"
	"fmt"
	"star/constant/str"
	"star/models"
	"testing"
	"time"
)

func Test_getFuncNewInstance(t *testing.T) {
	message := &models.PrivateMessage{
		Id:          111,
		SenderId:    222,
		RecipientId: 333,
		Content:     "hello world",
		Status:      false,
		SendTime:    time.Now(),
	}
	messageJson, err := json.Marshal(message)
	if err != nil {
		t.Error("json marshal error", err.Error())
		return
	}
	get := getFuncNewInstance(str.MessagePrivateMsg)
	pmessage := get()
	if err := json.Unmarshal(messageJson, pmessage); err != nil {
		t.Error("json unmarshal error", err.Error())
		return
	}
	fmt.Println(pmessage)
}
