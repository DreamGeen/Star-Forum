package service

import (
	"context"
	"star/app/comment/dao/mysql"
	"star/app/comment/dao/redis"
	"star/proto/comment/commentPb"
)

// StarComment 点赞评论
func (s *CommentService) StarComment(ctx context.Context, req *commentPb.StarCommentRequest, rsp *commentPb.StarCommentResponse) error {
	// redis中点赞
	if err := redis.IncrementCommentStar(req.CommentId); err != nil {
		rsp.Success = false
		rsp.Message = err.Error()
		return err
	}
	// mysql中点赞
	if err := mysql.UpdateStar(req.CommentId, 1); err != nil {
		rsp.Success = false
		rsp.Message = err.Error()
		return err
	}
	// 从redis中获取更新后的点赞数
	star, err := redis.GetCommentStar(req.CommentId)
	// 缓存未命中或已过期
	if err != nil || star == 0 {
		// 如果redis获取失败或返回0，则从MySQL中获取点赞数
		if dbStar, err := mysql.GetStar(req.CommentId); err != nil {
			rsp.Success = false
			rsp.Message = err.Error()
			return err
		} else {
			// 将从MySQL获取的点赞数更新回redis缓存
			_ = redis.SetCommentStar(req.CommentId, dbStar)
			rsp.Star = dbStar
		}
	} else {
		rsp.Star = star
	}

	rsp.Success = true
	rsp.Message = "点赞成功"
	return nil
}
