package jobs

import (
	"github.com/fadilmartias/firavel/config"
	"github.com/hibiken/asynq"
)

var (
	AsynqClient *asynq.Client
	AsynqServer *asynq.Server
)

func InitQueue(redisClient *config.RedisClient) {
	// Pakai redis options yang sama
	AsynqClient = asynq.NewClient(asynq.RedisClientOpt{
		Addr:     redisClient.GetClient().Options().Addr,
		Password: redisClient.GetClient().Options().Password,
		DB:       redisClient.GetClient().Options().DB,
	})

	AsynqServer = asynq.NewServer(
		asynq.RedisClientOpt{
			Addr:     redisClient.GetClient().Options().Addr,
			Password: redisClient.GetClient().Options().Password,
			DB:       redisClient.GetClient().Options().DB,
		},
		asynq.Config{
			Concurrency: 4,
			Queues: map[string]int{
				"critical": 6,
				"high":     4,
				"default":  3,
				"low":      1,
			},
		},
	)
}
