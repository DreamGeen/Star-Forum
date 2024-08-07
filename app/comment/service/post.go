package service

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"star/app/comment/dao/redis"
	"star/app/gateway/middleware/RabbitMQ"
	"star/models"
	"star/proto/comment/commentPb"
	"star/utils"
)

type CommentService struct{}

// PostComment 发布评论
func (s *CommentService) PostComment(ctx context.Context, req *commentPb.PostCommentRequest, rsp *commentPb.PostCommentResponse) error {
	comment := &models.Comment{
		PostId:      req.PostId,
		UserId:      req.UserId,
		Content:     req.Content,
		BeCommentId: req.BeCommentId,
	}

	// 尝试将评论存储到Redis
	if err := redis.CreateComment(comment); err != nil {
		// 如果存储到Redis失败，记录错误并返回
		utils.Logger.Error("存储评论到Redis失败", zap.Error(err))
		return err
	}

	// 使用RabbitMQ异步存储至MySQL数据库中
	// 生产者发布评论消息
	go func() {
		if err := RabbitMQ.PublishCommentEvent(comment); err != nil {
			if err := redis.Client.Del(ctx, fmt.Sprintf("comment:%d", comment.CommentId)).Err(); err != nil {
				utils.Logger.Error("删除Redis中评论缓存失败", zap.Error(err))
			}
			rsp.Success = false
			rsp.Message = err.Error()
		}
	}()

	rsp.Success = true
	rsp.Message = "评论发布成功"

	return nil
}
