package utils

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"sync"
)

// GatewayLogger 网关服务日志
var GatewayLogger *zap.Logger

var once sync.Once // 确保Logger只被初始化一次

// InitGatewayLogger 网关服务日志初始化
func InitGatewayLogger() error {
	once.Do(func() {
		encoder := getEncoder()

		// gateway.log记录全量日志
		logF, err := os.Create("D:\\Star-Forum\\Star-Forum\\app\\gateway\\logger\\log\\gateway.log")
		if err != nil {
			return
		}
		c1 := zapcore.NewCore(encoder, zapcore.AddSync(logF), zapcore.DebugLevel)

		// gateway.err.log记录ERROR级别的日志
		errF, err := os.Create("D:\\Star-Forum\\Star-Forum\\app\\gateway\\logger\\log\\gateway.err.log")
		if err != nil {
			return
		}
		c2 := zapcore.NewCore(encoder, zapcore.AddSync(errF), zap.ErrorLevel)

		core := zapcore.NewTee(c1, c2)
		GatewayLogger = zap.New(core, zap.AddCaller())
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
