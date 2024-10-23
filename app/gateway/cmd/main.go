package main

import (
	"context"
	"fmt"
	"github.com/go-micro/plugins/v4/registry/etcd"
	"go-micro.dev/v4/registry"
	"go-micro.dev/v4/web"
	"go.uber.org/zap"
	"star/app/constant/settings"
	"star/app/extra/tracing"
	"star/app/gateway/client"
	"star/app/gateway/router"
	"star/app/utils/logging"
)

func main() {
	//初始化微服务客户端
	client.Init()
	etcdReg := etcd.NewRegistry(
		registry.Addrs(fmt.Sprintf("%s:%d", settings.Conf.EtcdHost, settings.Conf.EtcdPort)))
	tp, err := tracing.SetTraceProvider("HttpService")
	if err != nil {
		logging.Logger.Error("set tracer error",
			zap.Error(err))
		return
	}
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			logging.Logger.Error("set tracer error",
				zap.Error(err))
			return
		}
	}()
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
