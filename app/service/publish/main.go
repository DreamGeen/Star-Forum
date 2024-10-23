package main

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"star/app/constant/settings"
	"star/app/constant/str"
	"star/app/extra/tracing"
	"star/app/utils/logging"
	"star/app/utils/snowflake"
	"star/proto/publish/publishPb"

	"github.com/go-micro/plugins/v4/registry/etcd"
	"go-micro.dev/v4"
	"go-micro.dev/v4/registry"
)

func main() {
	//雪花算法初始化
	if err := snowflake.Init(settings.Conf.SnowflakeId); err != nil {
		panic(err)
	}
	tp, err := tracing.SetTraceProvider(str.PublishService)
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

	publishSrvIns.New()
	//etcd注册件
	etcdReg := etcd.NewRegistry(
		registry.Addrs(fmt.Sprintf("%s:%d", settings.Conf.EtcdHost, settings.Conf.EtcdPort)),
	)
	//得到一个微服务实例
	microService := micro.NewService(
		micro.Name(str.PublishService), //服务名称
		micro.Version("v1"),            //服务版本
		micro.Registry(etcdReg),        //etcd注册件
	)
	//服务注册
	if err := publishPb.RegisterPublishServiceHandler(microService.Server(), publishSrvIns); err != nil {
		panic(err)
	}
	//服务启动
	if err := microService.Run(); err != nil {
		panic(err)
	}

}
