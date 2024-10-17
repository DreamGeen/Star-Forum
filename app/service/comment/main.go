package main

import (
	"fmt"
	"github.com/go-micro/plugins/v4/registry/etcd"
	"go-micro.dev/v4"
	"go-micro.dev/v4/registry"
	"log"
	"star/app/constant/settings"
	"star/app/constant/str"
	"star/app/gateway/client"
	"star/app/utils/snowflake"
	"star/proto/comment/commentPb"
)

func main() {
	//// 初始化zap
	//if err := logger.InitCommentLogger(); err != nil {
	//	log.Fatalf("初始化日志失败: %v", err)
	//}
	//
	//// 确保所有日志都被刷新
	//defer func() {
	//	if err := logger.CommentLogger.Sync(); err != nil {
	//		// 如果日志刷新失败，打印到标准错误输出
	//		_, _ = fmt.Fprintf(os.Stderr, "日志刷新失败: %v\n", err)
	//	}
	//}()

	//// 初始化配置
	//if err := settings.Init(); err != nil {
	//	log.Println("初始化配置失败", err)
	//}

	//雪花算法初始化
	if err := snowflake.Init(1); err != nil {
		log.Println("初始化雪花算法失败", err)
	}

	// etcd注册
	etcdReg := etcd.NewRegistry(
		registry.Addrs(fmt.Sprintf("%s:%d", settings.Conf.EtcdConfig.EtcdHost, settings.Conf.EtcdConfig.EtcdPort)),
	)

	// 初始化heartbeatStop channel
	heartbeatStop := make(chan struct{})
	// 停止发布心跳消息
	defer close(heartbeatStop)

	//feed.New()
	client.Init()

	// 创建服务
	microService := micro.NewService(
		micro.Name(str.CommentService),
		micro.Version("v1"),
		micro.Registry(etcdReg),
	)
	comment := &CommentService{}
	// 初始化服务
	// 级别会比 NewService 更高，作用一致，二选一即可
	// 后续代码运行期，初始化才有使用的必要
	//service.Init()

	// 注册服务
	if err := commentPb.RegisterCommentServiceHandler(microService.Server(), comment); err != nil {
		log.Println("注册服务失败", err)
	}
	// 运行服务
	if err := microService.Run(); err != nil {
		log.Println("运行服务失败", err)
	}
}
