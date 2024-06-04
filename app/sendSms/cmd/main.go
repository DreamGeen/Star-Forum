package main

import (
	"fmt"

	"github.com/go-micro/plugins/v4/registry/etcd"
	"go-micro.dev/v4"
	"go-micro.dev/v4/registry"

	redis "star/app/sendSms/dao"
	"star/app/sendSms/service"
	"star/proto/sendSms/sendSmsPb"
	"star/settings"
)

func main() {
	//初始化配置
	if err := settings.Init(); err != nil {
		panic(err)
	}
	//初始化redis
	if err := redis.Init(); err != nil {
		panic(err)
	}
	defer redis.Close()
	etcdReg := etcd.NewRegistry(
		registry.Addrs(fmt.Sprintf("%s:%d", settings.Conf.EtcdConfig.Host, settings.Conf.EtcdConfig.Port)),
	)
	//得到一个微服务实例
	microSevice := micro.NewService(
		micro.Name("SendSmsService"), //服务名称
		micro.Version("v1"),          //服务版本
		micro.Registry(etcdReg),      //etcd注册件
	)
	//初始化
	microSevice.Init()
	//服务注册
	if err := sendSmsPb.RegisterSendMsgHandler(microSevice.Server(), service.GetSendSmsSrv()); err != nil {
		panic(err)
	}
	//服务启动
	if err := microSevice.Run(); err != nil {
		panic(err)
	}
}
