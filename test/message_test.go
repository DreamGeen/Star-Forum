package test

import (
	"context"
	"star/app/gateway/client"
	"star/proto/message/messagePb"
	"testing"
)




func Test_SendSystemMessage(t *testing.T) {
	req := &messagePb.SendSystemMessageRequest{
		ManagerId:   1820019310731464704,
		RecipientId: 0,
		Type:        "all",
		Title:       "Test",
		Content:     "test send system message",
	}
	resp, err := client.SendSystemMessage(context.Background(), req)
	if err != nil {
		t.Error(err)
	}
	t.Log(resp)
}
