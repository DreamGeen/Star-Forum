package main

import (
	"fmt"
	"github.com/go-micro/plugins/v4/registry/etcd"
	"go-micro.dev/v4"
	"go-micro.dev/v4/registry"
	"go.uber.org/zap"
	"log"
	"os"
	"star/app/comment/dao/mysql"
	"star/app/comment/dao/redis"
	logger "star/app/comment/logger"
	"star/app/comment/rabbitMQ"
	commentService "star/app/comment/service"
	"star/proto/comment/commentPb"
	"star/settings"
	"star/utils"
	"time"
)

func main() {
	// 初始化zap
	if err := logger.InitCommentLogger(); err != nil {
		log.Fatalf("初始化日志失败: %v", err)
	}

	// 确保所有日志都被刷新
	defer func() {
		if err := logger.CommentLogger.Sync(); err != nil {
			// 如果日志刷新失败，打印到标准错误输出
			_, _ = fmt.Fprintf(os.Stderr, "日志刷新失败: %v\n", err)
		}
	}()

	// 初始化配置
	if err := settings.Init(); err != nil {
		logger.CommentLogger.Fatal("初始化配置失败", zap.Error(err))
	}

	// 初始化MySQL
	if err := mysql.Init(); err != nil {
		logger.CommentLogger.Fatal("初始化MySQL失败", zap.Error(err))
	}
	defer mysql.Close()

	// 初始化Redis
	if err := redis.Init(); err != nil {
		logger.CommentLogger.Fatal("初始化Redis失败", zap.Error(err))
	}
	defer redis.Close()

	//雪花算法初始化
	if err := utils.Init(1); err != nil {
		logger.CommentLogger.Fatal("初始化雪花算法失败", zap.Error(err))
	}

	// etcd注册
	etcdReg := etcd.NewRegistry(
		registry.Addrs(fmt.Sprintf("%s:%d", settings.Conf.EtcdConfig.EtcdHost, settings.Conf.EtcdConfig.EtcdPort)),
	)

	// 初始化RabbitMQ连接
	if err := rabbitMQ.ConnectToRabbitMQ(); err != nil {
		logger.CommentLogger.Fatal("初始化RabbitMQ连接失败", zap.Error(err))
	}
	defer rabbitMQ.Close()

	// 初始化heartbeatStop channel
	heartbeatStop := make(chan struct{})
	// 停止发布心跳消息
	defer close(heartbeatStop)

	// 发布心跳消息
	go rabbitMQ.StartHeartbeatTicker("comment_star", 5*time.Minute, heartbeatStop)
	go rabbitMQ.StartHeartbeatTicker("comment_delete", 5*time.Minute, heartbeatStop)
	//go rabbitMQ.StartHeartbeatTicker("comment_post", 5*time.Minute, heartbeatStop)

	// 启动RabbitMQ消费者
	rabbitMQ.ConsumeStarEvents()
	rabbitMQ.ConsumeDeleteCommentEvents()
	//RabbitMQ.ConsumeCommentEvents()

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
		logger.CommentLogger.Fatal("注册服务失败", zap.Error(err))
	}

	// 运行服务
	if err := service.Run(); err != nil {
		logger.CommentLogger.Fatal("运行服务失败", zap.Error(err))
	}
}
