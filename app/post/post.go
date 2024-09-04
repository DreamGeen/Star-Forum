package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"star/app/storage/cached"
	"star/app/storage/mysql"
	"star/constant/str"
	"star/proto/post/postPb"
)

type PostSrv struct {
}

var post *PostSrv

func (p *PostSrv) QueryPostExist(ctx context.Context, req *postPb.QueryPostExistRequest, resp *postPb.QueryPostExistResponse) error {
	key := fmt.Sprintf("QueryPostExist:%d", req.PostId)
	if _, err := cached.GetWithFunc(ctx, key, func(key string) (string, error) {
		return mysql.QueryPostExist(req.PostId)
	}); err != nil {
		if errors.Is(err, str.ErrPostNotExists) {
			return str.ErrPostNotExists
		}
		log.Println("query post exist err:", err)
		return str.ErrPostError
	}
	return nil
}
