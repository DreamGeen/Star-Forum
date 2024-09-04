package main

import (
	"fmt"
	"github.com/go-micro/plugins/v4/registry/etcd"
	"go-micro.dev/v4"
	"go-micro.dev/v4/registry"
	"log"
	"star/constant/settings"
	"star/constant/str"
	"star/proto/post/postPb"
	"star/utils"
)

func main() {
	////初始化配置
	//if err := settings.Init(); err != nil {
	//	panic(err)
	//}
	//雪花算法初始化
	if err := utils.Init(1); err != nil {
		panic(err)
	}

	//etcd注册件
	etcdReg := etcd.NewRegistry(
		registry.Addrs(fmt.Sprintf("%s:%d", settings.Conf.EtcdHost, settings.Conf.EtcdPort)),
	)
	//得到一个微服务实例
	microSevice := micro.NewService(
		micro.Name(str.PostService), //服务名称
		micro.Version("v1"),         //服务版本
		micro.Registry(etcdReg),     //etcd注册件
	)
	//初始化
	microSevice.Init()
	//服务注册
	if err := postPb.RegisterPostHandler(microSevice.Server(), post); err != nil {
		log.Println("failed to register post handler")
		panic(err)
	}
	//服务启动
	if err := microSevice.Run(); err != nil {
		panic(err)
	}

}
