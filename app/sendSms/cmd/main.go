package main

import (
	"fmt"
	"go.uber.org/zap"
	"log"
	"os"
	"star/app/gateway/middleware/RabbitMQ"
	"star/utils"

	"github.com/go-micro/plugins/v4/registry/etcd"
	"go-micro.dev/v4"
	"go-micro.dev/v4/registry"

	redis "star/app/sendSms/dao"
	"star/app/sendSms/service"
	"star/proto/sendSms/sendSmsPb"
	"star/settings"
)

func main() {
	// 初始化zap
	if err := utils.InitLogger(); err != nil {
		log.Fatalf("初始化日志失败: %v", err)
	}

	// 确保所有日志都被刷新
	defer func() {
		if err := utils.Logger.Sync(); err != nil {
			// 如果日志刷新失败，打印到标准错误输出
			_, _ = fmt.Fprintf(os.Stderr, "日志刷新失败: %v\n", err)
		}
	}()

	//初始化配置
	if err := settings.Init(); err != nil {
		utils.Logger.Fatal("初始化配置失败", zap.Error(err))
	}

	//初始化redis
	if err := redis.Init(); err != nil {
		utils.Logger.Fatal("初始化Redis失败", zap.Error(err))
	}
	defer redis.Close()

	//初始化etcd
	etcdReg := etcd.NewRegistry(
		registry.Addrs(fmt.Sprintf("%s:%d", settings.Conf.EtcdHost, settings.Conf.EtcdPort)),
	)

	// 初始化RabbitMQ连接
	if err := RabbitMQ.ConnectToRabbitMQ(); err != nil {
		utils.Logger.Fatal("初始化RabbitMQ连接失败", zap.Error(err))
	}
	defer RabbitMQ.Close()

	//得到一个微服务实例
	microSevice := micro.NewService(
		micro.Name("SendSmsService"), //服务名称
		micro.Version("v1"),          //服务版本
		micro.Registry(etcdReg),      //etcd注册件
	)
	//初始化
	// 级别会比 NewService 更高，二选一即可
	//microSevice.Init()

	//服务注册
	if err := sendSmsPb.RegisterSendMsgHandler(microSevice.Server(), service.GetSendSmsSrv()); err != nil {
		utils.Logger.Fatal("注册服务失败", zap.Error(err))
	}

	//服务启动
	if err := microSevice.Run(); err != nil {
		utils.Logger.Fatal("启动服务失败", zap.Error(err))
	}
}
