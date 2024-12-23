package redis

import (
	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"

	"loki/internal/config"
)

type Redis interface {
	Connection() *redis.Client
}

type redisClient struct {
	client *redis.Client
}

func NewRedisClient(cfg *config.Config) (Redis, error) {
	options, err := redis.ParseURL(cfg.RedisURI)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(options)

	if err = redisotel.InstrumentTracing(client); err != nil {
		return nil, err
	}
	if err = redisotel.InstrumentMetrics(client); err != nil {
		return nil, err
	}

	return &redisClient{client: client}, nil
}

func (r *redisClient) Connection() *redis.Client {
	return r.client
}
