package client

import (
	"context"
	"star/proto/user/userPb"
)

func LoginPassword(ctx context.Context, in *userPb.LSRequest) (*userPb.LoginResponse, error) {
	return userServiceClient.LoginPassword(ctx, in)

}

func LoginCaptcha(ctx context.Context, in *userPb.LSRequest) (*userPb.LoginResponse, error) {
	return userServiceClient.LoginCaptcha(ctx, in)
}

func Signup(ctx context.Context, in *userPb.LSRequest) (*userPb.EmptyLSResponse, error) {
	return userServiceClient.Signup(ctx, in)
}
