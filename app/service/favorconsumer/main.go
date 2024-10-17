package favorconsumer

import (
	"encoding/json"
	"github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
	str2 "star/app/constant/str"
	models2 "star/app/models"
	"star/app/storage/mysql"
	"star/app/utils/logging"
	"star/app/utils/rabbitmq"
	"sync"
	"time"
)

type collectMsg struct {
	count   int64
	userIds map[int64]struct{}
}

var conn *amqp091.Connection
var channel *amqp091.Channel

var likePostBuffer map[int64]int

var likeCommentBuffer map[int64]int

var collectPostBuffer map[int64]*collectMsg
var deleteCollect map[int64][]int64

var muLikePost sync.RWMutex
var muLikeComment sync.RWMutex
var muCollect sync.RWMutex

func failOnError(err error, msg string) {
	if err != nil {
		logging.Logger.Error(msg, zap.Error(err))
	}
}

func CloseMQ() {
	if err := conn.Close(); err != nil {
		logging.Logger.Error("close rabbitmq conn error", zap.Error(err))
		panic(err)
	}
	if err := channel.Close(); err != nil {
		logging.Logger.Error("close rabbitmq channel error", zap.Error(err))
		panic(err)
	}
}
func main() {
	var err error
	conn, err = amqp091.Dial(rabbitmq.ReturnRabbitmqUrl())
	failOnError(err, "Failed to connect to RabbitMQ")

	channel, err = conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer CloseMQ()

	go likePostConsumer()
	go likeCommentConsumer()
	go asyncUpdatePostLike()
	go asyncUpdateCommentLike()
	go collectPostConsumer()
	go asyncUpdatePostCollect()
}

func likePostConsumer() {
	delivery, err := channel.Consume(str2.LikePost, str2.Empty,
		false, false, false, false, nil)
	failOnError(err, "Failed to register a consumer")

	for msg := range delivery {
		var message *models2.Post
		if err := json.Unmarshal(msg.Body, &message); err != nil {
			logging.Logger.Error("unmarshal json fail")
			continue
		}
		muLikePost.Lock()
		likePostBuffer[message.PostId] += message.Star
		muLikePost.Unlock()
	}
}
func likeCommentConsumer() {
	delivery, err := channel.Consume(str2.LikeComment, str2.Empty,
		false, false, false, false, nil)
	failOnError(err, "Failed to register a consumer")

	for msg := range delivery {
		var comment *models2.Comment
		if err := json.Unmarshal(msg.Body, &comment); err != nil {
			logging.Logger.Error("unmarshal json fail")
			continue
		}
		muLikeComment.Lock()
		likeCommentBuffer[comment.PostId] += comment.Star
		muLikeComment.Unlock()
	}
}

func collectPostConsumer() {
	delivery, err := channel.Consume(str2.CollectPost, str2.Empty,
		false, false, false, false, nil)
	failOnError(err, "Failed to register a consumer")

	for msg := range delivery {
		var post *models2.Collect
		if err := json.Unmarshal(msg.Body, &post); err != nil {
			logging.Logger.Error("unmarshal json fail")
			continue
		}
		muCollect.Lock()
		collectBuffer := collectPostBuffer[post.PostId]
		collectBuffer.count += post.Collection
		if post.Collection > 0 {
			collectBuffer.userIds[post.UserId] = struct{}{}
		} else {
			if _, exists := collectBuffer.userIds[post.UserId]; exists {
				delete(collectBuffer.userIds, post.UserId)
			} else {
				deleteCollect[post.PostId] = append(deleteCollect[post.PostId], post.UserId)
			}
		}
		muCollect.Unlock()
	}
}
func asyncUpdatePostLike() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			muLikePost.Lock()
			for postId, count := range likePostBuffer {
				if err := mysql.UpdatePostLike(postId, count); err != nil {
					logging.Logger.Error("async feed update like fail", zap.Error(err), zap.Int64("post_id", postId), zap.Int("count", count))
				}
				delete(likePostBuffer, postId)
			}
			muLikePost.Unlock()
		}
	}
}
func asyncUpdateCommentLike() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			muLikeComment.Lock()
			for commentId, count := range likeCommentBuffer {
				if err := mysql.UpdateCommentLike(commentId, count); err != nil {
					logging.Logger.Error("async comment update like fail", zap.Error(err), zap.Int64("comment_id", commentId), zap.Int("count", count))
				}
				delete(likeCommentBuffer, commentId)
			}
			muLikeComment.Unlock()
		}
	}
}

func asyncUpdatePostCollect() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			muCollect.Lock()
			for postId, userCollect := range collectPostBuffer {
				if err := mysql.AddCollect(postId, userCollect); err != nil {
					logging.Logger.Error("async feed add collectMsg fail",
						zap.Error(err),
						zap.Int64("post_id", postId),
						zap.Int("count", userCollect.count))
				}
				delete(collectPostBuffer, postId)
			}
			for postId, userIds := range deleteCollect {
				if err := mysql.DeleteCollect(postId, userIds); err != nil {
					logging.Logger.Error("async feed delete collectMsg fail",
						zap.Error(err),
						zap.Int64("post_id", postId))
				}
				delete(deleteCollect, postId)
			}
			muCollect.Unlock()
		}
	}
}
