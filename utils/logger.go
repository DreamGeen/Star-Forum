package utils

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"star/constant/settings"
)

// Logger 创建一个全局的日志变量
var Logger *zap.Logger

// 初始化Logger
func init() {
	writeSyncer := getLogWriter(settings.Conf.FileName, settings.Conf.MaxSize, settings.Conf.MaxBackups, settings.Conf.MaxAge)
	encoder := getEncoder()
	var l = new(zapcore.Level)
	//将原配置中的字符串Lever反序列化为zap的lever级别并赋值给l
	err := l.UnmarshalText([]byte(settings.Conf.Level))
	if err != nil {
		panic(err)
	}
	core := zapcore.NewCore(encoder, writeSyncer, l)
	//将设置的配置传入，赋值给全局Logger变量
	Logger = zap.New(core, zap.AddCaller())
	return
}

// 设置编码器，即如何写入日志
func getEncoder() zapcore.Encoder {
	//返回zap默认编码结构体
	encoderConfig := zap.NewProductionEncoderConfig()

	//IS08601 UTC 时间格式("2006-01-02T15:04:05.000Z0700")
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	encoderConfig.TimeKey = "time"

	//将日志级别Lever序列化为全大写字符串
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	return zapcore.NewJSONEncoder(encoderConfig)
}

// 将日志配置信息导入
func getLogWriter(filename string, maxSize, maxBackup, maxAge int) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    maxSize,
		MaxBackups: maxBackup,
		MaxAge:     maxAge,
	}
	return zapcore.AddSync(lumberJackLogger)
}
