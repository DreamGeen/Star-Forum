package utils

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"sync"
)

// CommentLogger 评论服务日志
var CommentLogger *zap.Logger

var once sync.Once // 确保Logger只被初始化一次

// InitCommentLogger 评论服务日志初始化
func InitCommentLogger() error {
	once.Do(func() {
		encoder := getEncoder()

		// comment.log记录全量日志
		logF, err := os.Create("D:\\Star-Forum\\Star-Forum\\app\\comment\\logger\\log\\comment.log")
		if err != nil {
			return
		}
		c1 := zapcore.NewCore(encoder, zapcore.AddSync(logF), zapcore.DebugLevel)

		// comment.err.log记录ERROR级别的日志
		errF, err := os.Create("D:\\Star-Forum\\Star-Forum\\app\\comment\\logger\\log\\comment.err.log")
		if err != nil {
			return
		}
		c2 := zapcore.NewCore(encoder, zapcore.AddSync(errF), zap.ErrorLevel)

		core := zapcore.NewTee(c1, c2)
		CommentLogger = zap.New(core, zap.AddCaller())
	})
	return nil
}

// 设置日志的格式
func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}
