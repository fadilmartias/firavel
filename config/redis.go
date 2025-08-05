package config

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
	Prefix   string
}

type RedisClient struct {
	client *redis.Client
	prefix string
}

func LoadRedisConfig() RedisConfig {
	db, _ := strconv.Atoi(os.Getenv("REDIS_DB"))
	return RedisConfig{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       db,
		Prefix:   os.Getenv("REDIS_PREFIX"),
	}
}

func NewRedisClient() *RedisClient {
	redisConfig := LoadRedisConfig()
	client := redis.NewClient(&redis.Options{
		Addr:     redisConfig.Addr,
		Password: redisConfig.Password,
		DB:       redisConfig.DB,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}

	return &RedisClient{
		client: client,
		prefix: redisConfig.Prefix + ":",
	}
}

func (r *RedisClient) key(k string) string {
	return r.prefix + k
}

func (r *RedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.client.Set(ctx, r.key(key), value, expiration).Err()
}

func (r *RedisClient) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, r.key(key)).Result()
}

func (r *RedisClient) Del(ctx context.Context, key string) error {
	return r.client.Del(ctx, r.key(key)).Err()
}

func (r *RedisClient) GetClient() *redis.Client {
	return r.client
}

func (r *RedisClient) Close() error {
	return r.client.Close()
}

func (r *RedisClient) FlushDB() error {
	return r.client.FlushDB(context.Background()).Err()
}

func (r *RedisClient) FlushAll() error {
	return r.client.FlushAll(context.Background()).Err()
}

func (r *RedisClient) Keys(ctx context.Context, pattern string) ([]string, error) {
	return r.client.Keys(ctx, pattern).Result()
}

func (r *RedisClient) DeleteKeysByPrefix(ctx context.Context, prefix string) error {
	batchSize := 500 // jumlah key dihapus per batch
	keysToDelete := make([]string, 0, batchSize)

	iter := r.client.Scan(ctx, 0, prefix+"*", 0).Iterator()
	for iter.Next(ctx) {
		keysToDelete = append(keysToDelete, iter.Val())

		// Jika sudah mencapai batchSize, hapus dengan UNLINK biar non-blocking
		if len(keysToDelete) >= batchSize {
			if err := r.client.Unlink(ctx, keysToDelete...).Err(); err != nil {
				return err
			}
			keysToDelete = keysToDelete[:0]
		}
	}

	// Hapus sisa key kalau ada
	if len(keysToDelete) > 0 {
		if err := r.client.Unlink(ctx, keysToDelete...).Err(); err != nil {
			return err
		}
	}

	if err := iter.Err(); err != nil {
		return err
	}

	return nil
}
