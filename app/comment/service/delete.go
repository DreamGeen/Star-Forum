package service

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"star/app/comment/dao/mysql"
	"star/app/comment/dao/redis"
	"star/proto/comment/commentPb"
	"star/utils"
)

// DeleteComment 删除评论
func (s *CommentService) DeleteComment(ctx context.Context, req *commentPb.DeleteCommentRequest, rsp *commentPb.DeleteCommentResponse) error {
	if err := mysql.DeleteComment(req.CommentId); err != nil {
		rsp.Success = false
		rsp.Message = err.Error()
		return err
	}
	// 清除Redis缓存中与该评论相关的点赞数和回复数
	if err := redis.Client.Del(ctx, fmt.Sprintf("comment:star:%d", req.CommentId)).Err(); err != nil {
		utils.Logger.Error("删除Redis中点赞数缓存失败", zap.Error(err))
	}
	if err := redis.Client.Del(ctx, fmt.Sprintf("comment:reply:%d", req.CommentId)).Err(); err != nil {
		utils.Logger.Error("删除Redis中回复数缓存失败", zap.Error(err))
	}
	rsp.Success = true
	rsp.Message = "评论删除成功"
	return nil
}
