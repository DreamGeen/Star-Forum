package favorconsumer

import (
	"encoding/json"
	"github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
	"star/app/storage/mq"
	"star/app/storage/mysql"
	"star/constant/str"
	"star/models"
	"star/utils"
	"sync"
	"time"
)

var conn *amqp091.Connection
var channel *amqp091.Channel

var likePostBuffer map[int64]int
var likeCommentBuffer map[int64]int
var mu sync.RWMutex

func failOnError(err error, msg string) {
	if err != nil {
		utils.Logger.Error(msg, zap.Error(err))
	}
}

func CloseMQ() {
	if err := conn.Close(); err != nil {
		utils.Logger.Error("close rabbitmq conn error", zap.Error(err))
		panic(err)
	}
	if err := channel.Close(); err != nil {
		utils.Logger.Error("close rabbitmq channel error", zap.Error(err))
		panic(err)
	}
}
func main() {
	var err error
	conn, err = amqp091.Dial(mq.ReturnRabbitmqUrl())
	failOnError(err, "Failed to connect to RabbitMQ")

	channel, err = conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer CloseMQ()

	go likePostConsumer()
	go likeCommentConsumer()
	go asyncUpdatePostLike()
	go asyncUpdateCommentLike()
}

func likePostConsumer() {
	delivery, err := channel.Consume(str.LikePost, str.Empty,
		false, false, false, false, nil)
	failOnError(err, "Failed to register a consumer")

	for msg := range delivery {
		var message *models.Post
		if err := json.Unmarshal(msg.Body, &message); err != nil {
			utils.Logger.Error("unmarshal json fail")
			continue
		}
		mu.Lock()
		likePostBuffer[message.PostId] += message.Star
		mu.Unlock()
	}
}
func likeCommentConsumer() {
	delivery, err := channel.Consume(str.LikeComment, str.Empty,
		false, false, false, false, nil)
	failOnError(err, "Failed to register a consumer")

	for msg := range delivery {
		var comment *models.Comment
		if err := json.Unmarshal(msg.Body, &comment); err != nil {
			utils.Logger.Error("unmarshal json fail")
			continue
		}
		mu.Lock()
		likeCommentBuffer[comment.PostId] += comment.Star
		mu.Unlock()
	}
}
func asyncUpdatePostLike() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			mu.Lock()
			for postId, count := range likePostBuffer {
				if err := mysql.UpdatePostLike(postId, count); err != nil {
					utils.Logger.Error("async post update like fail", zap.Error(err), zap.Int64("post_id", postId), zap.Int("count", count))
				}
				delete(likePostBuffer, postId)
			}
			mu.Unlock()
		}
	}
}
func asyncUpdateCommentLike() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			mu.Lock()
			for commentId, count := range likeCommentBuffer {
				if err := mysql.UpdateCommentLike(commentId, count); err != nil {
					utils.Logger.Error("async comment update like fail", zap.Error(err), zap.Int64("comment_id", commentId), zap.Int("count", count))
				}
				delete(likeCommentBuffer, commentId)
			}
			mu.Unlock()
		}
	}
}
