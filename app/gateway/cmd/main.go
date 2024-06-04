package main

import (
	"fmt"
	"star/app/gateway/client"

	"github.com/go-micro/plugins/v4/registry/etcd"
	"go-micro.dev/v4/registry"
	"go-micro.dev/v4/web"

	"star/app/gateway/router"
	"star/settings"
)

func main() {
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
	_ = webService.Init()
	_ = webService.Run()
}
