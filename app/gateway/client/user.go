package client

import (
	"context"
	"star/proto/user/userPb"
)

func LoginPassword(ctx context.Context, in *userPb.LSRequest) (*userPb.LoginResponse, error) {
	return userService.LoginPassword(ctx, in)

}

func LoginCaptcha(ctx context.Context, in *userPb.LSRequest) (*userPb.LoginResponse, error) {
	return userService.LoginCaptcha(ctx, in)
}

func Signup(ctx context.Context, in *userPb.LSRequest) (*userPb.EmptyLSResponse, error) {
	return userService.Signup(ctx, in)
}
