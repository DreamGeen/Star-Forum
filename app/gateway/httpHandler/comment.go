package httpHandler

import (
	"star/app/constant/str"
	"star/app/extra/tracing"
	"star/app/gateway/client"
	"star/app/gateway/models"
	"star/app/utils/logging"
	"star/proto/comment/commentPb"
	"strconv"

	"go.uber.org/zap"

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
	_, span := tracing.Tracer.Start(c.Request.Context(), "PostCommentHandler")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger := logging.LogServiceWithTrace(span, "GateWay.PostComment")

	// 参数校验
	p := new(models.PostComment)
	if err := c.ShouldBindJSON(p); err != nil {
		logger.Error("feed comment error,invalid param",
			zap.Error(err))
		str.Response(c, str.ErrInvalidParam, nil)
		return
	}
	// 发布处理
	req := &commentPb.PostCommentRequest{
		PostId:      p.PostId,
		UserId:      p.UserId,
		Content:     p.Content,
		BeCommentId: p.BeCommentId,
	}
	resp, err := client.PostComment(c.Request.Context(), req)
	if err != nil {
		logger.Error("feed comment error",
			zap.Error(err),
			zap.Int64("userId", req.UserId),
			zap.String("content", req.Content),
			zap.Int64("BeCommentId", req.BeCommentId))
		str.Response(c, err, nil)
		return
	}

	// 成功响应
	str.Response(c, nil, map[string]interface{}{
		"comment": resp.Content,
	})
}

func DeleteComment(c *gin.Context) {
	// 测试样例：127.0.0.1:9090/comment/1
	_, span := tracing.Tracer.Start(c.Request.Context(), "DeleteCommentHandler")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger := logging.LogServiceWithTrace(span, "GateWay.DeleteComment")

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		logger.Error("delete comment error,invalid param",
			zap.Error(err))
		str.Response(c, str.ErrInvalidParam, nil)
		return
	}
	req := &commentPb.DeleteCommentRequest{CommentId: id}
	_, err = client.DeleteComment(c, req)
	if err != nil {
		logger.Error("delete comment error",
			zap.Error(err),
			zap.Int64("commentId", id))
		str.Response(c, err, nil)
		return
	}
	// 成功响应
	str.Response(c, nil, nil)
}

func GetComments(c *gin.Context) {
	// 测试样例：127.0.0.1:9090/comments?postId=1
	_, span := tracing.Tracer.Start(c.Request.Context(), "GetCommentsHandler")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger := logging.LogServiceWithTrace(span, "GateWay.GetComments")

	postId, err := strconv.ParseInt(c.Query("postId"), 10, 64)
	if err != nil {
		logger.Error("get comment error,invalid param",
			zap.Error(err))
		str.Response(c, str.ErrInvalidParam, nil)
		return
	}
	req := &commentPb.GetCommentsRequest{PostId: postId}
	resp, err := client.GetComments(c.Request.Context(), req)
	if err != nil {
		logger.Error("get comment error",
			zap.Error(err))
		str.Response(c, err, nil)
		return
	}

	// 成功响应
	str.Response(c, nil, map[string]interface{}{
		"comments": resp.Comments,
	})
}
