package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rabbitmq/amqp091-go"
	redis2 "github.com/redis/go-redis/v9"
	"github.com/robfig/cron/v3"
	"go-micro.dev/v4"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"star/app/constant/str"
	"star/app/extra/tracing"
	"star/app/models"
	"star/app/storage/cached"
	"star/app/storage/mysql"
	"star/app/storage/redis"
	"star/app/utils/logging"
	"star/app/utils/rabbitmq"
	"star/app/utils/snowflake"
	"star/proto/message/messagePb"
	"star/proto/user/userPb"
	"sync"
	"time"
)

type MessageSrv struct {
}

var userService userPb.UserService
var conn *amqp091.Connection
var channel *amqp091.Channel
var manager = NewManager()

func failOnError(err error, msg string) {
	if err != nil {
		logging.Logger.Error(msg, zap.Error(err))
	}
}

func CloseMQ() {
	if err := conn.Close(); err != nil {
		logging.Logger.Error("message service close rabbitmq conn error",
			zap.Error(err))
		panic(err)
	}
	if err := channel.Close(); err != nil {
		logging.Logger.Error("message service close rabbitmq channel error",
			zap.Error(err))
		panic(err)
	}
}

func (m *MessageSrv) New() {
	//连接消息队列
	var err error
	conn, err = amqp091.Dial(rabbitmq.ReturnRabbitmqUrl())
	failOnError(err, "message service failed to connect to RabbitMQ")

	channel, err = conn.Channel()
	failOnError(err, "message service failed to open a channel")

	err = channel.ExchangeDeclare(str.MessageExchange,
		"topic",
		false, false, false, false,
		nil)
	failOnError(err, "message service failed to declare an exchange")
	//声明队列
	_, err = channel.QueueDeclare(str.MessageLike,
		false, false, false, false,
		nil)
	failOnError(err, "message service failed to declare a like queue")
	_, err = channel.QueueDeclare(str.MessageReply,
		false, false, false, false,
		nil)
	failOnError(err, "message service failed to declare a reply queue")
	_, err = channel.QueueDeclare(str.MessageSystem,
		false, false, false, false,
		nil)
	failOnError(err, "message service failed to declare a system queue")
	_, err = channel.QueueDeclare(str.MessagePrivateMsg,
		false, false, false, false,
		nil)
	failOnError(err, "message service failed to declare a private_msg queue")
	_, err = channel.QueueDeclare(str.MessageMention,
		false, false, false, false,
		nil)
	failOnError(err, "message service failed to declare a mention queue")

	//绑定队列
	// 绑定点赞消息
	err = channel.QueueBind(str.MessageLike, str.RoutMessageLike, str.MessageExchange, false, nil)
	failOnError(err, "message service failed to bind like queue")

	// 绑定@提及消息
	err = channel.QueueBind(str.MessageMention, str.RoutMention, str.MessageExchange, false, nil)
	failOnError(err, "message service failed to bind mention queue")

	// 绑定回复消息
	err = channel.QueueBind(str.MessageReply, str.RoutMention, str.MessageExchange, false, nil)
	failOnError(err, "message service failed to bind reply queue")

	// 绑定系统通知
	err = channel.QueueBind(str.MessageSystem, str.RoutSystem, str.MessageExchange, false, nil)
	failOnError(err, "message service failed to bind system queue")

	// 绑定私信消息
	err = channel.QueueBind(str.MessagePrivateMsg, str.RoutPrivateMsg, str.MessageExchange, false, nil)
	failOnError(err, "message service failed to bind private message queue")

	//创建一个用户微服务客户端
	userMicroService := micro.NewService(micro.Name(str.UserServiceClient))
	userService = userPb.NewUserService(str.UserService, userMicroService.Client())

	cronRunner := cron.New()
	cronRunner.AddFunc("0 2 * * * ?", removeMessage)

	cronRunner.Start()

	go manager.Run()

}

// ListMessageCount 获取未读消息数量
func (m *MessageSrv) ListMessageCount(ctx context.Context, req *messagePb.ListMessageCountRequest, resp *messagePb.ListMessageCountResponse) error {
	ctx, span := tracing.Tracer.Start(ctx, "ListMessageCountService")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger := logging.LogServiceWithTrace(span, "MessageService.ListMessageCount")

	//获取countsStr
	key := fmt.Sprintf("ListMessageCount:+%d", req.UserId)
	countsStr, err := cached.GetWithFunc(ctx, key, func(key string) (string, error) {
		counts, err := mysql.ListMessageCount(req.UserId)
		if err != nil {
			logger.Error("ListMessageCount service error",
				zap.Error(err))
			logging.SetSpanError(span, err)
			return "", str.ErrMessageError
		}
		countsJson, err := json.Marshal(counts)
		if err != nil {
			logger.Error("json marshal counts error:",
				zap.Error(err),
				zap.Any("counts", counts))
			logging.SetSpanError(span, err)
			return "", err
		}
		return string(countsJson), nil

	})
	if err != nil {
		logger.Error("ListMessageCount service error",
			zap.Error(err),
			zap.Int64("userId", req.UserId))
		logging.SetSpanError(span, err)
		return err
	}
	//解析countsStr
	counts := new(models.Counts)
	err = json.Unmarshal([]byte(countsStr), counts)
	if err != nil {
		logger.Error("json unmarshal counts error",
			zap.Error(err),
			zap.Int64("userId", req.UserId))
		logging.SetSpanError(span, err)
		return str.ErrMessageError
	}
	resp.Count = &messagePb.Counts{
		PrivateMsgCount: counts.PrivateMsgCount,
		MentionCount:    counts.MentionCount,
		LikeCount:       counts.LikeCount,
		ReplyCount:      counts.ReplyCount,
		SystemCount:     counts.SystemCount,
		TotalCount:      counts.TotalCount,
	}
	return nil
}

// SendSystemMessage 发送系统消息
func (m *MessageSrv) SendSystemMessage(ctx context.Context, req *messagePb.SendSystemMessageRequest, resp *messagePb.SendSystemMessageResponse) error {
	ctx, span := tracing.Tracer.Start(ctx, "SendSystemMessageService")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger := logging.LogServiceWithTrace(span, "MessageService.SendSystemMessage")

	message := &models.SystemMessage{
		Id:          snowflake.GetID(),
		RecipientId: req.RecipientId,
		ManagerId:   req.ManagerId,
		Type:        req.Type,
		Title:       req.Title,
		Content:     req.Content,
		Status:      false,
		PublishTime: time.Now().UTC(),
	}
	body, err := json.Marshal(message)
	if err != nil {
		logger.Error("json marshal system message error",
			zap.Error(err),
			zap.Any("message", message))
		logging.SetSpanError(span, err)
		return str.ErrMessageError
	}
	header := rabbitmq.InjectAMQPHeaders(ctx)
	err = channel.PublishWithContext(
		ctx,
		str.MessageExchange,
		str.RoutSystem,
		false,
		false,
		amqp091.Publishing{
			DeliveryMode: amqp091.Persistent,
			ContentType:  "text/plain",
			Body:         body,
			Headers:      header,
		})
	if err != nil {
		logger.Error("send system message error",
			zap.Error(err),
			zap.Any("message", message))
		logging.SetSpanError(span, err)
		return str.ErrMessageError
	}
	return nil
}

// SendPrivateMessage 私信
func (m *MessageSrv) SendPrivateMessage(ctx context.Context, req *messagePb.SendPrivateMessageRequest, resp *messagePb.SendPrivateMessageResponse) error {
	ctx, span := tracing.Tracer.Start(ctx, "SendPrivateMessageService")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger := logging.LogServiceWithTrace(span, "MessageService.SendPrivateMessage")

	message := &models.PrivateMessage{
		Id:            snowflake.GetID(),
		SenderId:      req.SenderId,
		RecipientId:   req.RecipientId,
		Content:       req.Content,
		Status:        false,
		SendTime:      time.Now().UTC(),
		PrivateChatId: req.PrivateChatId,
	}
	if err := redis.SaveMessage(ctx, message); err != nil {
		logger.Error("redis save private message error",
			zap.Error(err),
			zap.Any("message", message))
		logging.SetSpanError(span, err)
		return str.ErrMessageError
	}
	body, err := json.Marshal(message)
	if err != nil {
		logger.Error("json marshal message error",
			zap.Error(err),
			zap.Any("message", message))
		logging.SetSpanError(span, err)
		return err
	}
	header := rabbitmq.InjectAMQPHeaders(ctx)
	err = channel.PublishWithContext(
		ctx,
		str.MessageExchange,
		str.RoutPrivateMsg,
		false,
		false,
		amqp091.Publishing{
			DeliveryMode: amqp091.Persistent,
			ContentType:  "text/plain",
			Body:         body,
			Headers:      header,
		},
	)
	if err != nil {
		logger.Error("publish message error",
			zap.Error(err),
			zap.Any("message", message))
		logging.SetSpanError(span, err)
		return err
	}
	client, ok := manager.GetClient(req.RecipientId)
	if ok {
		client.send <- message
	}
	return nil
}

// SendRemindMessage 发送提醒消息
func (m *MessageSrv) SendRemindMessage(ctx context.Context, req *messagePb.SendRemindMessageRequest, resp *messagePb.SendRemindMessageResponse) error {
	ctx, span := tracing.Tracer.Start(ctx, "SendRemindMessageService")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger := logging.LogServiceWithTrace(span, "MessageService.SendRemindMessage")

	var err error
	switch req.RemindType {
	case "like":
		err = addRemindMessage(ctx, req, str.RoutMessageLike, span, logger)
	case "reply":
		err = addRemindMessage(ctx, req, str.RoutReply, span, logger)
	case "mention":
		err = addRemindMessage(ctx, req, str.RoutMention, span, logger)
	}
	if err != nil {
		logger.Error("add message error",
			zap.Error(err),
			zap.Int64("senderId", req.SenderId),
			zap.Int64("recipientId", req.RecipientId),
			zap.String("content", req.Content),
			zap.String("url", req.Url),
			zap.String("sourceType", req.SourceType),
			zap.Int64("sourceId", req.SourceId))
		logging.SetSpanError(span, err)
		return str.ErrMessageError
	}
	return nil
}

func addRemindMessage(ctx context.Context, req *messagePb.SendRemindMessageRequest, routingKey string, span trace.Span, logger *zap.Logger) error {
	message := &models.RemindMessage{
		Id:          snowflake.GetID(),
		SourceId:    req.SourceId,
		SourceType:  req.SourceType,
		SenderId:    req.SenderId,
		RecipientId: req.RecipientId,
		Content:     req.Content,
		Url:         req.Url,
		Status:      false,
		RemindTime:  time.Now().UTC(),
	}
	body, err := json.Marshal(message)
	if err != nil {
		logger.Error("json marshal message error",
			zap.Error(err),
			zap.Any("message", message))
		logging.SetSpanError(span, err)
		return err
	}
	header := rabbitmq.InjectAMQPHeaders(ctx)
	err = channel.PublishWithContext(
		ctx,
		str.MessageExchange,
		routingKey,
		false,
		false,
		amqp091.Publishing{
			DeliveryMode: amqp091.Persistent,
			ContentType:  "text/plain",
			Body:         body,
			Headers:      header,
		})
	if err != nil {
		logger.Error("publish message error",
			zap.Error(err),
			zap.Any("message", message))
		logging.SetSpanError(span, err)
		return err
	}
	return nil
}

// GetChatList 获取会话列表
func (m *MessageSrv) GetChatList(ctx context.Context, req *messagePb.GetChatListRequest, resp *messagePb.GetChatListResponse) error {
	ctx, span := tracing.Tracer.Start(ctx, "GetChatListService")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger := logging.LogServiceWithTrace(span, "MessageService.GetChatList")

	key := fmt.Sprintf("chatList:%d", req.UserId)
	val, err := redis.Client.Get(ctx, key).Result()
	if err != nil {
		if !errors.Is(err, redis2.Nil) {
			logger.Error("redis get chatList error",
				zap.Error(err),
				zap.String("key", key),
				zap.Int64("UserId", req.UserId))
			logging.SetSpanError(span, err)
			return str.ErrMessageError
		}
		//为空则去mysql里查询
		list, err := mysql.GetChatList(req.UserId)
		if err != nil {
			logger.Error("mysql get chatList error",
				zap.Error(err),
				zap.Int64("UserId", req.UserId))
			logging.SetSpanError(span, err)
			return str.ErrMessageError
		}
		if len(list) == 0 {
			return nil
		}
		privateChatList, err := convertChatListToPB(ctx, req.UserId, list, logger)
		if err != nil {
			logger.Error("convertChatListToPB error",
				zap.Error(err),
				zap.Int64("UserId", req.UserId))
			logging.SetSpanError(span, err)
			return str.ErrMessageError
		}
		resp.PrivateChatList = privateChatList
		privateListJosn, err := json.Marshal(privateChatList)
		if err != nil {
			logger.Error("json marshal chatList error",
				zap.Error(err),
				zap.Any("list", list))
			logging.SetSpanError(span, err)
			return str.ErrMessageError
		}
		redis.Client.Set(ctx, key, string(privateListJosn), 24*time.Hour)
		return nil
	}
	var list []*messagePb.PrivateChat
	if err := json.Unmarshal([]byte(val), &list); err != nil {
		logger.Error("json unmarshal message error",
			zap.Error(err),
			zap.Any("value", val))
		logging.SetSpanError(span, err)
		return str.ErrMessageError
	}
	resp.PrivateChatList = list
	return nil
}

func getSenderId(chat *models.PrivateChat, recipientId int64) int64 {
	if chat.User1Id == recipientId {
		return chat.User1Id
	}
	return chat.User2Id
}

func convertChatListToPB(ctx context.Context, recipientId int64, list []*models.PrivateChat, logger *zap.Logger) ([]*messagePb.PrivateChat, error) {
	plist := make([]*messagePb.PrivateChat, len(list))
	var wg sync.WaitGroup
	var chatChan = make(chan struct {
		index int
		pchat *messagePb.PrivateChat
	}, len(list))
	for i, chat := range list {
		wg.Add(1)
		go func(i int, chat *models.PrivateChat) {
			defer wg.Done()
			senderResp, err := userService.GetUserInfo(ctx, &userPb.GetUserInfoRequest{
				UserId: getSenderId(chat, recipientId),
			})
			if err != nil {
				logger.Error("get sender user info error",
					zap.Error(err),
					zap.Int64("userId", recipientId))
				return
			}
			sender := senderResp.User
			pchat := &messagePb.PrivateChat{
				UserId:      sender.UserId,
				UserName:    sender.UserName,
				Img:         sender.Img,
				LastMsg:     chat.LastMsgContent,
				LastMsgTime: chat.LastSendTime.Format(str.ParseTimeFormat),
			}
			chatChan <- struct {
				index int
				pchat *messagePb.PrivateChat
			}{index: i, pchat: pchat}
		}(i, chat)
	}
	go func() {
		wg.Wait()
		close(chatChan)
	}()
	for chat := range chatChan {
		plist[chat.index] = chat.pchat
	}
	return plist, nil
}

// LoadMessage 加载消息
func (m *MessageSrv) LoadMessage(ctx context.Context, req *messagePb.LoadMessageRequest, resp *messagePb.LoadMessageResponse) error {
	ctx, span := tracing.Tracer.Start(ctx, "LoadMessageService")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger := logging.LogServiceWithTrace(span, "MessageService.LoadMessage")

	lastMsgTime, err := time.Parse(str.ParseTimeFormat, req.LastMsgTime)
	if err != nil {
		logger.Error("load last message time error",
			zap.Error(err),
			zap.Int64("userId", req.RecipientId))
		logging.SetSpanError(span, err)
		return str.ErrMessageError
	}
	senderReq, err := userService.GetUserInfo(ctx, &userPb.GetUserInfoRequest{
		UserId: req.SenderId,
	})
	if err != nil {
		logger.Error("get sender user info error",
			zap.Error(err),
			zap.Int64("senderId", req.SenderId))
		logging.SetSpanError(span, err)
		return str.ErrMessageError
	}
	sender := senderReq.User
	messages, err := redis.LoadMessage(ctx, req.SenderId, req.RecipientId, lastMsgTime, str.DefaultLoadMessageNumber)
	if err != nil {
		logger.Error("load last message error",
			zap.Error(err),
			zap.Int64("senderId", req.SenderId),
			zap.Int64("recipientId", req.RecipientId))
		logging.SetSpanError(span, err)
		return str.ErrMessageError
	}
	if len(messages) > 0 {
		//redis找到了
		resp.PrivateMessages = convertMessagesToPB(messages, sender)
		return nil
	}
	//redis没有找到
	messages, err = mysql.LoadMessage(req.PrivateChatId, lastMsgTime, 2*str.DefaultLoadMessageNumber)
	if err != nil {
		logger.Error("load last message error",
			zap.Error(err),
			zap.Int64("privateChatId", req.PrivateChatId))
		logging.SetSpanError(span, err)
		return str.ErrMessageError
	}
	resp.PrivateMessages = convertMessagesToPB(messages, sender)
	//将message异步存入redis
	go func() {
		if err := redis.BitchSaveMessage(ctx, messages); err != nil {
			logger.Error("async redis save private message error",
				zap.Error(err),
				zap.Any("message", messages))
		}
	}()
	return nil
}

func convertMessagesToPB(messages []*models.PrivateMessage, sender *userPb.User) []*messagePb.PrivateMessage {
	pmessages := make([]*messagePb.PrivateMessage, len(messages))
	for i, message := range messages {
		pmessage := &messagePb.PrivateMessage{
			SenderId:    message.SenderId,
			SenderName:  sender.UserName,
			SenderImg:   sender.Img,
			RecipientId: message.RecipientId,
			Content:     message.Content,
			Status:      message.Status,
			SendTime:    message.SendTime.Format(str.ParseTimeFormat),
		}
		pmessages[i] = pmessage
	}
	return pmessages
}

func removeMessage() {
	ctx, span := tracing.Tracer.Start(context.Background(), "removeMessageService")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger := logging.LogServiceWithTrace(span, "MessageService.removeMessage")

	goroutineLimiter := make(chan struct{}, 15)
	chats, err := mysql.GetAllPrivateChat()
	if err != nil {
		logger.Error("mysql get all private chat error",
			zap.Error(err))
		logging.SetSpanError(span, err)
		return
	}
	var wg sync.WaitGroup
	for _, chat := range chats {
		wg.Add(1)
		goroutineLimiter <- struct{}{}
		go func() {
			defer func() {
				wg.Done()
				<-goroutineLimiter
			}()
			length := redis.GetMessageLength(ctx, chat.User1Id, chat.User2Id)
			if length > 200 {
				err := redis.RemoveMessage(ctx, chat.User1Id, chat.User2Id, 0, 100)
				if err != nil {
					logger.Error("remove message error",
						zap.Error(err),
						zap.Int64("chatId", chat.Id),
						zap.Int64("user1Id", chat.User1Id),
						zap.Int64("user2Id", chat.User2Id))
				}
			}

		}()
	}
	wg.Wait()
	close(goroutineLimiter)
}

func (m *MessageSrv) SendMessage(ctx context.Context, req *messagePb.SendMessageRequest, resp *messagePb.SendMessageResponse) error {

	return nil
}
