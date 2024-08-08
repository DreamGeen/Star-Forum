package main

import (
	"fmt"
	"github.com/go-micro/plugins/v4/registry/etcd"
	"go-micro.dev/v4"
	"go-micro.dev/v4/registry"
	"go.uber.org/zap"
	"log"
	"os"
	redis "star/app/sendSms/dao"
	logger "star/app/sendSms/logger"
	"star/app/sendSms/service"
	"star/proto/sendSms/sendSmsPb"
	"star/settings"
)

func main() {
	// 初始化zap
	if err := logger.InitSendSmsLogger(); err != nil {
		log.Fatalf("初始化日志失败: %v", err)
	}

	// 确保所有日志都被刷新
	defer func() {
		if err := logger.SendSmsLogger.Sync(); err != nil {
			// 如果日志刷新失败，打印到标准错误输出
			_, _ = fmt.Fprintf(os.Stderr, "日志刷新失败: %v\n", err)
		}
	}()

	//初始化配置
	if err := settings.Init(); err != nil {
		logger.SendSmsLogger.Fatal("初始化配置失败", zap.Error(err))
	}

	//初始化redis
	if err := redis.Init(); err != nil {
		logger.SendSmsLogger.Fatal("初始化Redis失败", zap.Error(err))
	}
	defer redis.Close()

	//初始化etcd
	etcdReg := etcd.NewRegistry(
		registry.Addrs(fmt.Sprintf("%s:%d", settings.Conf.EtcdHost, settings.Conf.EtcdPort)),
	)

	//得到一个微服务实例
	microService := micro.NewService(
		micro.Name("SendSmsService"), //服务名称
		micro.Version("v1"),          //服务版本
		micro.Registry(etcdReg),      //etcd注册件
	)
	//初始化
	// 级别会比 NewService 更高，二选一即可
	//microService.Init()

	//服务注册
	if err := sendSmsPb.RegisterSendMsgHandler(microService.Server(), service.GetSendSmsSrv()); err != nil {
		logger.SendSmsLogger.Fatal("注册服务失败", zap.Error(err))
	}

	//服务启动
	if err := microService.Run(); err != nil {
		logger.SendSmsLogger.Fatal("启动服务失败", zap.Error(err))
	}
}
