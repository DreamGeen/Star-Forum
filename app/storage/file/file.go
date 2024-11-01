package file

import (
	"context"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"go.uber.org/zap"
	"mime/multipart"
	"star/app/constant/settings"
	"star/app/constant/str"
	"star/app/utils/logging"
)

// UploadToQiNiu 封装上传图片或文件到七牛云然后返回状态和url
func UploadToQiNiu(ctx context.Context, file multipart.File, fileSize int64) (string, error) {
	var AccessKey = settings.Conf.AccessKey
	var SecretKey = settings.Conf.SecretKey
	var Bucket = settings.Conf.Bucket
	var ImgUrl = settings.Conf.QiniuServer
	putPlicy := storage.PutPolicy{
		Scope: Bucket,
	}
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
	err := formUploader.PutWithoutKey(ctx, &ret, upToken, file, fileSize, &putExtra)
	if err != nil {
		logging.Logger.Error("upload image error:",
			zap.Error(err))
		return "", str.ErrUpload
	}
	url := ImgUrl + "/" + ret.Key
	return url, nil
}
