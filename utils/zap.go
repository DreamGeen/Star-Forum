package utils

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"sync"
)

// Logger 在其它包中使用utils.Logger即可调用
var Logger *zap.Logger
var once sync.Once // 确保Logger只被初始化一次

func InitLogger() error {
	var err error
	once.Do(func() {
		encoder := getEncoder()

		// star.log记录全量日志
		logF, err := os.Create("D:\\Star-Forum\\Star-Forum\\log\\star.log")
		if err != nil {
			return
		}
		c1 := zapcore.NewCore(encoder, zapcore.AddSync(logF), zapcore.DebugLevel)

		// star.err.log记录ERROR级别的日志
		errF, err := os.Create("D:\\Star-Forum\\Star-Forum\\log\\star.err.log")
		if err != nil {
			return
		}
		c2 := zapcore.NewCore(encoder, zapcore.AddSync(errF), zap.ErrorLevel)

		core := zapcore.NewTee(c1, c2)
		Logger = zap.New(core, zap.AddCaller())
	})
	return err
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}
