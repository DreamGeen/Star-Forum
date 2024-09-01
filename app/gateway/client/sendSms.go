package client

import (
	"context"
	"star/proto/sendSms/sendSmsPb"
)

func HandleSendSms(ctx context.Context, in *sendSmsPb.SendRequest) (*sendSmsPb.EmptySendResponse, error) {
	return sendSmsService.HandleSendSms(ctx, in)
}
