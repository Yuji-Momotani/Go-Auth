package cache

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisClient interface {
	Set(ctx context.Context, key string, val string, expipration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
}

type redisClient struct {
	client *redis.Client
}

func NewRedisClient() RedisClient {
	addr := fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"))
	redis := redis.NewClient(&redis.Options{Addr: addr})

	// テスト：終わったら削除
	if err := redis.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("failed redis.NewClient: %s", err)
	}

	client := redisClient{
		client: redis,
	}

	return &client
}

func (c *redisClient) Set(
	ctx context.Context,
	key string,
	val string,
	expipration time.Duration,
) error {
	err := c.client.Set(ctx, key, val, expipration).Err()
	if err != nil {
		return err
	}

	return nil
}

func (c *redisClient) Get(
	ctx context.Context,
	key string,
) (string, error) {
	cmd := c.client.Get(ctx, key)
	val, err := cmd.Result()
	if err != nil {
		return "", cmd.Err()
	}

	return val, nil
}
