package httpHandler

import (
	"star/app/constant/str"
	"star/app/extra/tracing"
	"star/app/gateway/client"
	"star/app/gateway/models"
	"star/app/utils/logging"
	"star/app/utils/request"
	"star/proto/publish/publishPb"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func PreUploadVideosHandler(c *gin.Context) {
	_, span := tracing.Tracer.Start(c.Request.Context(), "PreUploadVideos")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger := logging.LogServiceWithTrace(span, "GateWay.PreUploadVideos")

	preUpload := new(models.PreUpload)
	if err := c.ShouldBind(preUpload); err != nil {
		logger.Error("preUploadVideos error because invalid param",
			zap.Error(err))
		str.Response(c, str.ErrInvalidParam, nil)
		return
	}
	userId, err := request.GetUserId(c)
	if err != nil {
		logger.Error("get user id error",
			zap.Error(err))
		str.Response(c, err, nil)
		return
	}
	resp, err := client.PreUploadVideos(c.Request.Context(), &publishPb.PreUploadVideosRequest{
		ActorId:  userId,
		Chunks:   preUpload.Chunks,
		FileName: preUpload.FileName,
	})
	if err != nil {
		logger.Error("pre upload videos error",
			zap.Error(err),
			zap.Int64("userId", userId))
		str.Response(c, err, nil)
		return
	}
	str.Response(c, nil, map[string]interface{}{
		"uploadId": resp.UploadId,
	})
}

func CreatePostHandler(c *gin.Context) {
	_, span := tracing.Tracer.Start(c.Request.Context(), "CreatePostHandler")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger := logging.LogServiceWithTrace(span, "GateWay.CreatePost")

	userId, err := request.GetUserId(c)
	if err != nil {
		logger.Warn("user not log in,but want to create feed")
		str.Response(c, err, nil)
		return
	}
	p := new(models.CreatePost)
	err = c.ShouldBindJSON(p)
	if err != nil {
		logger.Error("create feed error,invalid param",
			zap.Error(err),
			zap.Int64("userId", userId))
		str.Response(c, str.ErrInvalidParam, nil)
		return
	}
	_, err = client.CreatePost(c.Request.Context(), &publishPb.CreatePostRequest{
		UserId:      userId,
		IsScan:      p.IsScan,
		Content:     p.Content,
		CommunityId: p.CommunityId,
	})
	if err != nil {
		logger.Error("create feed  service error",
			zap.Error(err),
			zap.Int64("userId", userId),
			zap.Int64("communityId", p.CommunityId),
			zap.String("content", p.Content),
			zap.Bool("isScan", p.IsScan))
		str.Response(c, err, nil)
		return
	}
	str.Response(c, nil, nil)
	return
}
