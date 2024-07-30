package httpHandler

import (
	"net/http"
	"star/app/gateway/client"
	"star/proto/comment/commentPb"
	"strconv"

	"github.com/gin-gonic/gin"
)

func PostComment(c *gin.Context) {
	var req commentPb.PostCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp, err := client.PostComment(c, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func DeleteComment(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
		return
	}
	req := &commentPb.DeleteCommentRequest{CommentId: id}
	resp, err := client.DeleteComment(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func GetComments(c *gin.Context) {
	postId, err := strconv.ParseInt(c.Query("postId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}
	page, err := strconv.ParseInt(c.Query("page"), 10, 64)
	if err != nil {
		page = 1
	}
	pageSize, err := strconv.ParseInt(c.Query("pageSize"), 10, 64)
	if err != nil {
		pageSize = 10
	}
	req := &commentPb.GetCommentsRequest{PostId: postId, Page: page, PageSize: pageSize}
	resp, err := client.GetComments(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func StarComment(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
		return
	}
	req := &commentPb.StarCommentRequest{CommentId: id}
	resp, err := client.StarComment(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}
