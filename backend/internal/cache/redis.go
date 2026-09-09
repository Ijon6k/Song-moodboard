package cache

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache struct {
	Client *redis.Client
}

func Connect(addr, password string) (*Cache, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var err error
	for i := 0; i < 15; i++ {
		err = rdb.Ping(ctx).Err()
		if err == nil {
			break
		}
		log.Printf("Waiting for Redis to be ready... (%d/15): %v", i+1, err)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	log.Println("Connected to Redis successfully.")
	return &Cache{Client: rdb}, nil
}

func (c *Cache) Get(ctx context.Context, key string) (string, error) {
	return c.Client.Get(ctx, key).Result()
}

func (c *Cache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return c.Client.Set(ctx, key, value, expiration).Err()
}

func (c *Cache) Del(ctx context.Context, keys ...string) error {
	return c.Client.Del(ctx, keys...).Err()
}

func (c *Cache) BlacklistToken(ctx context.Context, token string, exp time.Duration) error {
	return c.Set(ctx, "bl:"+token, "revoked", exp)
}

func (c *Cache) IsTokenBlacklisted(ctx context.Context, token string) bool {
	val, err := c.Get(ctx, "bl:"+token)
	return err == nil && val == "revoked"
}
