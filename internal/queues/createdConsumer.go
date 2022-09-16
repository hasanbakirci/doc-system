package queues

import (
	"context"
	"github.com/hasanbakirci/doc-system/pkg/redisClient"
	"github.com/hasanbakirci/doc-system/pkg/response"
	log "github.com/sirupsen/logrus"
)

type CreatedConsumer struct {
	redisClient *redisClient.RedisClient
}

func NewCreatedConsumer(redis *redisClient.RedisClient) CreatedConsumer {
	return CreatedConsumer{redisClient: redis}
}

func (redis CreatedConsumer) Consume(channel string) {
	subs := redis.redisClient.Subscribe(channel)
	for {
		msg, err := subs.ReceiveMessage(context.Background())
		if err != nil {
			response.Panic(400, err.Error())
		}
		redis.redisClient.Set("doc-system:created-log", msg.Payload)
		log.Info(msg.Channel, msg.Payload)
	}
}
