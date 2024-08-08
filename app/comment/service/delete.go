package service

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"star/app/comment/dao/mysql"
	"star/app/comment/dao/redis"
	"star/app/gateway/middleware/RabbitMQ"
	"star/proto/comment/commentPb"
	"star/utils"
)

// DeleteComment 删除评论
func (s *CommentService) DeleteComment(ctx context.Context, req *commentPb.DeleteCommentRequest, rsp *commentPb.DeleteCommentResponse) error {
	// 检查评论是否存在
	if err := mysql.CheckComment(req.CommentId); err != nil {
		return err
	}

	// 清除对应评论的Redis缓存
	if err := redis.Client.Del(ctx, fmt.Sprintf("comment:star:%d", req.CommentId)).Err(); err != nil {
		utils.Logger.Error("删除Redis中点赞数缓存失败", zap.Error(err))
	}

	// 使用RabbitMQ异步在MySQL数据库中删除
	// 生产者发布删除消息
	if err := RabbitMQ.PublishDeleteEvent(req.CommentId); err != nil {
		return err
	}
	rsp.Success = true
	rsp.Message = "评论删除成功"
	return nil
}
