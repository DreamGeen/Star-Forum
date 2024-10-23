package main

import (
	"context"
	"fmt"
	"github.com/go-micro/plugins/v4/registry/etcd"
	"go-micro.dev/v4"
	"go-micro.dev/v4/registry"
	"go.uber.org/zap"
	"log"
	"star/app/constant/settings"
	"star/app/constant/str"
	"star/app/extra/tracing"
	"star/app/utils/logging"
	"star/app/utils/snowflake"
	"star/proto/comment/commentPb"
)

func main() {

	//雪花算法初始化
	if err := snowflake.Init(settings.Conf.SnowflakeId); err != nil {
		log.Println("初始化雪花算法失败", err)
	}
	tp, err := tracing.SetTraceProvider(str.CommentService)
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
	commentSrvIns.New()
	// etcd注册
	etcdReg := etcd.NewRegistry(
		registry.Addrs(fmt.Sprintf("%s:%d", settings.Conf.EtcdConfig.EtcdHost, settings.Conf.EtcdConfig.EtcdPort)),
	)

	// 初始化heartbeatStop channel
	heartbeatStop := make(chan struct{})
	// 停止发布心跳消息
	defer close(heartbeatStop)

	// 创建服务
	microService := micro.NewService(
		micro.Name(str.CommentService),
		micro.Version("v1"),
		micro.Registry(etcdReg),
	)

	// 初始化服务
	// 级别会比 NewService 更高，作用一致，二选一即可
	// 后续代码运行期，初始化才有使用的必要
	//service.Init()

	// 注册服务
	if err := commentPb.RegisterCommentServiceHandler(microService.Server(), commentSrvIns); err != nil {
		log.Println("注册服务失败", err)
	}
	// 运行服务
	if err := microService.Run(); err != nil {
		log.Println("运行服务失败", err)
	}
}
