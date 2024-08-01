package main

import (
	"fmt"
	"github.com/go-micro/plugins/v4/registry/etcd"
	"go-micro.dev/v4"
	"go-micro.dev/v4/logger"
	"go-micro.dev/v4/registry"
	"star/utils"

	"star/app/comment/dao/mysql"
	"star/app/comment/dao/redis"
	commentService "star/app/comment/service"
	"star/proto/comment/commentPb"
	"star/settings"
)

func main() {
	// 初始化配置
	if err := settings.Init(); err != nil {
		logger.Fatal(err)
	}

	// 初始化MySQL
	if err := mysql.Init(); err != nil {
		logger.Fatal(err)
	}
	defer mysql.Close()

	// 初始化Redis
	if err := redis.Init(); err != nil {
		logger.Fatal(err)
	}
	defer redis.Close()

	//雪花算法初始化
	if err := utils.Init(1); err != nil {
		panic(err)
	}

	// etcd注册
	etcdReg := etcd.NewRegistry(
		registry.Addrs(fmt.Sprintf("%s:%d", settings.Conf.EtcdConfig.EtcdHost, settings.Conf.EtcdConfig.EtcdPort)),
	)

	// 创建服务
	service := micro.NewService(
		micro.Name("CommentService"),
		micro.Version("v1"),
		micro.Registry(etcdReg),
	)

	// 初始化服务
	// 级别会比 NewService 更高，作用一致，二选一即可
	// 后续代码运行期，初始化才有使用的必要
	//service.Init()

	// 注册服务
	if err := commentPb.RegisterCommentServiceHandler(service.Server(), new(commentService.CommentService)); err != nil {
		logger.Fatal(err)
	}

	// 运行服务
	if err := service.Run(); err != nil {
		logger.Fatal(err)
	}
}
