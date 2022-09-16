package listener

import (
	"github.com/hasanbakirci/doc-system/internal/config"
	"github.com/hasanbakirci/doc-system/internal/queues"
	"github.com/hasanbakirci/doc-system/pkg/redisClient"
)

type listener struct {
	redisClient     *redisClient.RedisClient
	createdConsumer queues.CreatedConsumer
}

func NewListener(settings config.Configuration) listener {

	redis := redisClient.NewRedisClient(settings.RedisSettings.Uri)

	createdConsumer := queues.NewCreatedConsumer(redis)

	return listener{
		redisClient:     redis,
		createdConsumer: createdConsumer,
	}
}

func (receiver listener) Start() {
	go receiver.createdConsumer.Consume("doc-system")
}
