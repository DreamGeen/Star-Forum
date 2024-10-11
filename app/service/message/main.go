package main

import (
	"fmt"
	"github.com/go-micro/plugins/v4/registry/etcd"
	"go-micro.dev/v4"
	"go-micro.dev/v4/registry"
	"star/app/constant/settings"
	"star/app/constant/str"
	"star/app/utils/snowflake"
	"star/proto/message/messagePb"
)

func main() {
	//雪花算法初始化
	if err := snowflake.Init(1); err != nil {
		panic(err)
	}
	defer CloseMQ()
	message := new(MessageSrv)
	message.New()
	//etcd注册件
	etcdReg := etcd.NewRegistry(
		registry.Addrs(fmt.Sprintf("%s:%d", settings.Conf.EtcdHost, settings.Conf.EtcdPort)),
	)
	//得到一个微服务实例
	microSevice := micro.NewService(
		micro.Name(str.MessageService), //服务名称
		micro.Version("v1"),            //服务版本
		micro.Registry(etcdReg),        //etcd注册件
	)
	//服务注册
	if err := messagePb.RegisterMessageServiceHandler(microSevice.Server(), message); err != nil {
		panic(err)
	}
	//服务启动
	if err := microSevice.Run(); err != nil {
		panic(err)
	}

}
