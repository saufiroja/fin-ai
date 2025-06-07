package redis

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/saufiroja/fin-ai/config"
	logging "github.com/saufiroja/fin-ai/pkg/loggings"
)

type Redis interface {
	Set(key string, value interface{}) error
	Get(key string) (string, error)
	Del(key string) error
	Close() error
}

type RedisClient struct {
	client  *redis.Client
	context context.Context
}

var (
	instance *RedisClient
	once     sync.Once
)

func NewRedisClient(conf *config.AppConfig, log logging.Logger) *RedisClient {
	once.Do(func() {
		ctx := context.Background()

		rdb := redis.NewClient(&redis.Options{
			Addr:            conf.Redis.URL,
			Password:        "",
			DB:              0,
			PoolSize:        10,
			MinIdleConns:    3,
			DialTimeout:     5 * time.Second,
			ReadTimeout:     3 * time.Second,
			WriteTimeout:    3 * time.Second,
			PoolTimeout:     4 * time.Second,
			ConnMaxIdleTime: 10 * time.Minute,
		})

		if err := rdb.Ping(ctx).Err(); err != nil {
			log.LogError(fmt.Sprintf("failed to connect to Redis at %s: %v", conf.Redis.URL, err))
			return
		}

		log.LogInfo(fmt.Sprintf("successfully connected to Redis at %s", conf.Redis.URL))

		instance = &RedisClient{
			client:  rdb,
			context: ctx,
		}
	})

	return instance
}

func (r *RedisClient) Set(key string, value interface{}) error {
	_, err := r.client.Set(r.context, key, value, 0).Result()
	return err
}

func (r *RedisClient) Get(key string) (string, error) {
	value, err := r.client.Get(r.context, key).Result()
	if err != nil && err != redis.Nil {
		return "", err
	}
	return value, nil
}

func (r *RedisClient) Del(key string) error {
	_, err := r.client.Del(r.context, key).Result()
	return err
}

func (r *RedisClient) Close() error {
	return r.client.Close()
}
