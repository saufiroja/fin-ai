package redis

import (
	"context"
	"fmt"
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

func NewRedisClient(conf *config.AppConfig, logging logging.Logger) *RedisClient {
	ctx := context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr:     conf.Redis.URL, // Alamat Redis, misalnya "localhost:6379" atau "redis-server:6379"
		Password: "",             // Gunakan jika Redis pakai password
		DB:       0,              // Database yang digunakan (0–15)

		// Connection pool
		PoolSize:     10, // Jumlah maksimum koneksi dalam pool (disesuaikan dengan beban aplikasi)
		MinIdleConns: 3,  // Koneksi idle minimum yang dijaga tetap terbuka

		// Timeout
		DialTimeout:  5 * time.Second, // Timeout saat koneksi dibuat
		ReadTimeout:  3 * time.Second, // Timeout untuk operasi baca
		WriteTimeout: 3 * time.Second, // Timeout untuk operasi tulis

		// Health check
		PoolTimeout:     4 * time.Second,  // Timeout tunggu koneksi dari pool
		ConnMaxIdleTime: 10 * time.Minute, // Waktu maksimum koneksi idle sebelum ditutup
	})

	// Test connection
	if err := rdb.Ping(ctx).Err(); err != nil {
		logging.LogError(fmt.Sprintf("failed to connect to Redis at %s: %v", conf.Redis.URL, err))
		return nil
	}

	logging.LogInfo(fmt.Sprintf("successfully connected to Redis at %s", conf.Redis.URL))

	return &RedisClient{
		client:  rdb,
		context: ctx,
	}
}

func (r *RedisClient) Set(key string, value interface{}) error {
	_, err := r.client.Set(r.context, key, value, 0).Result()
	if err != nil {
		return err
	}
	return nil
}

func (r *RedisClient) Get(key string) (string, error) {
	value, err := r.client.Get(r.context, key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", nil // Key does not exist
		}
		return "", err
	}
	return value, nil
}

func (r *RedisClient) Del(key string) error {
	_, err := r.client.Del(r.context, key).Result()
	if err != nil {
		return err
	}
	return nil
}

func (r *RedisClient) Close() error {
	if err := r.client.Close(); err != nil {
		return fmt.Errorf("failed to close Redis client: %w", err)
	}
	return nil
}
