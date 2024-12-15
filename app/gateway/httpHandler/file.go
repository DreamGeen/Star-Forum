package httpHandler

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"star/app/constant/str"
	"star/app/extra/tracing"
	"star/app/storage/file"
	"star/app/utils/logging"
)

func FileUploadHandler(c *gin.Context) {
	_, span := tracing.Tracer.Start(c.Request.Context(), "FileUploadHandler")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger := logging.LogServiceWithTrace(span, "GateWay.FileUpload")

	f, err := c.FormFile("file")
	if err != nil {
		str.Response(c, str.ErrInvalidParam, nil)
		return
	}
	url, err := file.UploadToQiNiu(c.Request.Context(), str.DirImg, str.UploadMarkImg, f, logger)
	if err != nil {
		logger.Error("upload file error",
			zap.Error(err))
		str.Response(c, str.ErrUpload, nil)
		return
	}
	str.Response(c, nil, map[string]interface{}{
		"data": url,
	})
}
