package client

import (
	"context"
	"star/proto/message/messagePb"
)

func ListMessageCount(ctx context.Context, req *messagePb.ListMessageCountRequest) (*messagePb.ListMessageCountResponse, error) {
	return messageService.ListMessageCount(ctx, req)
}

func SendSystemMessage(ctx context.Context, req *messagePb.SendSystemMessageRequest) (*messagePb.SendSystemMessageResponse, error) {
	return messageService.SendSystemMessage(ctx, req)
}

func SendPrivateMessage(ctx context.Context, req *messagePb.SendPrivateMessageRequest) (*messagePb.SendPrivateMessageResponse, error) {
	return messageService.SendPrivateMessage(ctx, req)
}

func SendRemindMessage(ctx context.Context, req *messagePb.SendRemindMessageRequest) (*messagePb.SendRemindMessageResponse, error) {
	return messageService.SendRemindMessage(ctx, req)
}

func GetChatList(ctx context.Context, req *messagePb.GetChatListRequest) (*messagePb.GetChatListResponse, error) {
	return messageService.GetChatList(ctx, req)
}

func LoadMessage(ctx context.Context, req *messagePb.LoadMessageRequest) (*messagePb.LoadMessageResponse, error) {
	return messageService.LoadMessage(ctx, req)
}
