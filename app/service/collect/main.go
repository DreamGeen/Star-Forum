package collect

import (
	"context"
	"fmt"
	"star/app/constant/settings"
	"star/app/constant/str"
	"star/app/extra/tracing"
	"star/app/utils/logging"
	"star/app/utils/snowflake"
	"star/proto/collect/collectPb"

	"github.com/go-micro/plugins/v4/registry/etcd"
	"go-micro.dev/v4"
	"go-micro.dev/v4/registry"
	"go.uber.org/zap"
)

func main() {
	//雪花算法初始化
	if err := snowflake.Init(settings.Conf.SnowflakeId); err != nil {
		panic(err)
	}
	defer CloseMQ()
	tp, err := tracing.SetTraceProvider(str.CollectService)
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
	collectSrvIns.New()
	//etcd注册件
	etcdReg := etcd.NewRegistry(
		registry.Addrs(fmt.Sprintf("%s:%d", settings.Conf.EtcdHost, settings.Conf.EtcdPort)),
	)
	//得到一个微服务实例
	microService := micro.NewService(
		micro.Name(str.CollectService), //服务名称
		micro.Version("v1"),            //服务版本
		micro.Registry(etcdReg),        //etcd注册件
	)
	//服务注册
	if err := collectPb.RegisterCollectServiceHandler(microService.Server(), collectSrvIns); err != nil {
		panic(err)
	}
	//服务启动
	if err := microService.Run(); err != nil {
		panic(err)
	}
}
