package httpHandler

import (
	"go.uber.org/zap"
	"log"
	"star/app/gateway/client"
	logger "star/app/gateway/logger"
	"star/app/gateway/models"
	"star/constant/str"

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
		logger.GatewayLogger.Error("参数错误", zap.Error(err))
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
		log.Println("评论发布失败", err)
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
		log.Println("参数错误", err)
		str.Response(c, str.ErrInvalidParam, str.Empty, nil)
		return
	}
	req := &commentPb.DeleteCommentRequest{CommentId: id}
	_, err = client.DeleteComment(c, req)
	if err != nil {
		log.Println("评论删除失败", err)
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
		log.Println("参数错误", err)
		str.Response(c, str.ErrInvalidParam, str.Empty, nil)
		return
	}
	req := &commentPb.GetCommentsRequest{PostId: postId}
	resp, err := client.GetComments(c, req)
	if err != nil {
		log.Println("评论获取失败", err)
		str.Response(c, err, str.Empty, nil)
		return
	}
	// 成功响应
	str.Response(c, nil, "comments", resp.Comments)
}

func StarComment(c *gin.Context) {
	// 测试样例：127.0.0.1:9090/comment/star/1
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		log.Println("参数错误", err)
		str.Response(c, str.ErrInvalidParam, str.Empty, nil)
		return
	}
	req := &commentPb.StarCommentRequest{CommentId: id}
	resp, err := client.StarComment(c, req)
	if err != nil {
		log.Println("点赞评论失败", err)
		str.Response(c, err, str.Empty, nil)
		return
	}
	// 成功响应
	str.Response(c, nil, "star", resp.Star)
}
