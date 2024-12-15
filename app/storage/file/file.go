package file

import (
	"context"
	"github.com/google/uuid"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"go.uber.org/zap"
	"mime/multipart"
	"star/app/constant/settings"
	"star/app/constant/str"
	"strings"
	"time"
)

// UploadToQiNiu 封装上传图片或文件到七牛云然后返回状态和url
func UploadToQiNiu(ctx context.Context, dir string, uploadMark string, file *multipart.FileHeader, logger *zap.Logger) (string, error) {
	var AccessKey = settings.Conf.AccessKey
	var SecretKey = settings.Conf.SecretKey
	var Bucket = settings.Conf.Bucket
	var ImgUrl = settings.Conf.QiniuServer
	putPlicy := storage.PutPolicy{
		Scope: Bucket,
	}
	src, err := file.Open()
	if err != nil {
		return str.Empty, err
	}
	defer src.Close()
	mac := qbox.NewMac(AccessKey, SecretKey)
	upToken := putPlicy.UploadToken(mac)
	cfg := storage.Config{
		Zone:          &storage.ZoneHuadong,
		UseCdnDomains: false,
		UseHTTPS:      false,
	}
	putExtra := storage.PutExtra{}
	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}
	//生成新的文件名
	suffix := strings.Split(file.Filename, ".")[1]
	fileName := uuid.New().String() + uploadMark + str.Dot + suffix
	//生成完整上传路径
	fullDir := dir + time.Now().Format(str.DirTimeParse) + str.Backslashes
	key := fullDir + fileName
	//上传
	err = formUploader.Put(ctx, &ret, upToken, key, src, file.Size, &putExtra)
	if err != nil {
		logger.Error("upload image error:",
			zap.Error(err))
		return "", str.ErrUpload
	}
	url := ImgUrl + "/" + ret.Key
	return url, nil
}
