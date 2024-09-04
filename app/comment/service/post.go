package service

import (
	"context"
	"star/app/comment/dao/mysql"
	"star/models"
	"star/proto/comment/commentPb"
	"star/proto/post/postPb"
	"sync"
)

type CommentService struct {
	commentPb.CommentService
}

var (
	commentSrvIns *CommentService
	postService   postPb.PostService
	once          sync.Once
)

func GetCommentSrv() *CommentService {
	once.Do(func() {
		commentSrvIns = &CommentService{}
	})
	return commentSrvIns
}

//func (s *CommentService) New() {
//	postMicroService := micro.NewService(micro.Name(str.PostServiceClient))
//	postService = postPb.NewPostService(str.PostService, postMicroService.Client())
//}

// PostComment 发布评论
func (s *CommentService) PostComment(ctx context.Context, req *commentPb.PostCommentRequest, rsp *commentPb.PostCommentResponse) error {
	comment := &models.Comment{
		PostId:      req.PostId,
		UserId:      req.UserId,
		Content:     req.Content,
		BeCommentId: req.BeCommentId,
	}
	// 存储评论
	if err := mysql.CreateComment(comment); err != nil {
		return err
	}
	rsp.Content = comment.Content

	return nil
}
