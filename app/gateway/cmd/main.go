package main

import (
	"fmt"
	"github.com/go-micro/plugins/v4/registry/etcd"
	"go-micro.dev/v4/registry"
	"go-micro.dev/v4/web"
	"log"
	"os"
	"star/app/gateway/client"

	logger "star/app/gateway/logger"
	"star/app/gateway/router"
	"star/settings"
)

func main() {
	// 初始化zap
	if err := logger.InitGatewayLogger(); err != nil {
		log.Fatalf("初始化日志失败: %v", err)
	}

	// 确保所有日志都被刷新
	defer func() {
		if err := logger.GatewayLogger.Sync(); err != nil {
			// 如果日志刷新失败，打印到标准错误输出
			_, _ = fmt.Fprintf(os.Stderr, "日志刷新失败: %v\n", err)
		}
	}()

	// 初始化配置
	if err := settings.Init(); err != nil {
		fmt.Println(err)
		return
	}
	//初始化微服务客户端
	client.Init()
	etcdReg := etcd.NewRegistry(
		registry.Addrs(fmt.Sprintf("%s:%d", settings.Conf.EtcdHost, settings.Conf.EtcdPort)))
	//得到一个web服务实例
	webService := web.NewService(
		web.Name("HttpService"), //服务名称
		web.Address(fmt.Sprintf("%s:%d", settings.Conf.HttpHost, settings.Conf.HttpPort)),
		web.Registry(etcdReg),       // etcd注册件
		web.Handler(router.Setup()), // 路由
		web.Metadata(map[string]string{"protocol": "http"}),
	)
	//初始化并运行web服务
	_ = webService.Init()

	// 如果出错，检查端口是否被占用
	if err := webService.Run(); err != nil {
		fmt.Println(err)
	}
}
