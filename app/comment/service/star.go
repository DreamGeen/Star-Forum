package service

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"star/app/comment/dao/mysql"
	"star/app/comment/dao/redis"
	RabbitMQ "star/app/gateway/middleware/RabbitMQ"
	"star/proto/comment/commentPb"
	"star/utils"
)

// StarComment 点赞评论
func (s *CommentService) StarComment(ctx context.Context, req *commentPb.StarCommentRequest, rsp *commentPb.StarCommentResponse) error {
	// 尝试从Redis中获取点赞数
	star, err := redis.GetCommentStar(req.CommentId)
	// 缓存未命中或已过期
	if err != nil || star == 0 {
		// 如果Redis获取失败或返回0，则从MySQL中获取点赞数
		if dbStar, err := mysql.GetStar(req.CommentId); err != nil {
			rsp.Success = false
			rsp.Message = err.Error()
			return err
		} else {
			// 将从MySQL获取的点赞数更新回Redis缓存
			err = redis.SetCommentStar(req.CommentId, dbStar)
			if err != nil {
				utils.Logger.Error("更新点赞数至Redis失败", zap.Error(err))
			}
			rsp.Star = dbStar
		}
	} else {
		rsp.Star = star
	}

	// Redis中点赞
	if err := redis.IncrementCommentStar(req.CommentId); err != nil {
		utils.Logger.Error("Redis中点赞失败", zap.Error(err))
		rsp.Success = false
		rsp.Message = err.Error()
		return err
	}

	// 使用RabbitMQ异步存储至MySQL数据库中
	// 生产者发布点赞消息
	go func() {
		if err := RabbitMQ.PublishStarEvent(req.CommentId); err != nil {
			if err := redis.Client.Del(ctx, fmt.Sprintf("comment:star:%d", req.CommentId)).Err(); err != nil {
				utils.Logger.Error("删除Redis中点赞数缓存失败", zap.Error(err))
			}
			rsp.Success = false
			rsp.Message = err.Error()
		}
	}()

	rsp.Success = true
	rsp.Message = "点赞成功"
	rsp.Star = rsp.Star + 1
	return nil
}
