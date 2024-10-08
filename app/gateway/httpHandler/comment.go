package httpHandler

import (
	"go.uber.org/zap"
	"star/app/gateway/client"
	"star/app/gateway/models"
	"star/constant/str"
	"star/utils"

	"star/proto/comment/commentPb"
	"strconv"

	"github.com/gin-gonic/gin"
)

func PostComment(c *gin.Context) {
	// 测试样例：127.0.0.1:9090/comment
	//{
	//	"postId":1,
	//	"userId":123,
	//	"content":"家人们",
	//	"beCommentId":1819704700270809088
	//}
	// 参数校验
	p := new(models.PostComment)
	if err := c.ShouldBindJSON(p); err != nil {
		utils.Logger.Error("post comment error,invalid param",
			zap.Error(err))
		str.Response(c, str.ErrInvalidParam, str.Empty, nil)
		return
	}
	// 发布处理
	req := &commentPb.PostCommentRequest{
		PostId:      p.PostId,
		UserId:      p.UserId,
		Content:     p.Content,
		BeCommentId: p.BeCommentId,
	}
	resp, err := client.PostComment(c, req)
	if err != nil {
		utils.Logger.Error("post comment error",
			zap.Error(err),
			zap.Int64("userId", req.UserId),
			zap.String("content", req.Content),
			zap.Int64("BeCommentId", req.BeCommentId))
		str.Response(c, err, str.Empty, nil)
		return
	}
	// 成功响应
	str.Response(c, nil, "comment", resp.Content)
}

func DeleteComment(c *gin.Context) {
	// 测试样例：127.0.0.1:9090/comment/1
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.Logger.Error("delete comment error,invalid param",
			zap.Error(err))
		str.Response(c, str.ErrInvalidParam, str.Empty, nil)
		return
	}
	req := &commentPb.DeleteCommentRequest{CommentId: id}
	_, err = client.DeleteComment(c, req)
	if err != nil {
		utils.Logger.Error("delete comment error",
			zap.Error(err),
			zap.Int64("commentId", id))
		str.Response(c, err, str.Empty, nil)
		return
	}
	// 成功响应
	str.Response(c, nil, str.Empty, nil)
}

func GetComments(c *gin.Context) {
	// 测试样例：127.0.0.1:9090/comments?postId=1
	postId, err := strconv.ParseInt(c.Query("postId"), 10, 64)
	if err != nil {
		utils.Logger.Error("get comment error,invalid param",
			zap.Error(err))
		str.Response(c, str.ErrInvalidParam, str.Empty, nil)
		return
	}
	req := &commentPb.GetCommentsRequest{PostId: postId}
	resp, err := client.GetComments(c, req)
	if err != nil {
		utils.Logger.Error("get comment error",
			zap.Error(err))
		str.Response(c, err, str.Empty, nil)
		return
	}
	// 成功响应
	str.Response(c, nil, "comments", resp.Comments)
}
