package test

import (
	"context"
	"go-micro.dev/v4"
	"star/constant/str"
	"star/proto/message/messagePb"
	"testing"
)

var messageService messagePb.MessageService

func Init() {
	messageMicroService := micro.NewService(micro.Name(str.MessageServiceClient))
	messageService = messagePb.NewMessageService(str.MessageService, messageMicroService.Client())
}

func Test_SendSystemMessage(t *testing.T) {
	Init()
	req := &messagePb.SendSystemMessageRequest{
		ManagerId:   1820019310731464704,
		RecipientId: 0,
		Type:        "all",
		Title:       "Test",
		Content:     "test send system message",
	}
	resp, err := messageService.SendSystemMessage(context.Background(), req)
	if err != nil {
		t.Error(err)
	}
	t.Log(resp)
}
