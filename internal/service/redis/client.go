package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type Client interface {
	Set(ctx context.Context, key string, value any, expiration time.Duration) error
	SetNX(ctx context.Context, key string, value any, expiration time.Duration) (bool, error)
	Get(ctx context.Context, key string) (string, error)
	Del(ctx context.Context, key string) error
	Ping(ctx context.Context) error
	Close() error
	VerifyCaptcha(ctx context.Context, key, answer string) (string, error)
}

type client struct {
	rdb *redis.Client
}

func NewClient(url string) (Client, error) {
	opts, err := redis.ParseURL(url)
	if err != nil {
		return nil, err
	}

	rdb := redis.NewClient(opts)

	return &client{rdb: rdb}, nil
}

func (c *client) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	return c.rdb.Set(ctx, key, value, expiration).Err()
}

func (c *client) SetNX(ctx context.Context, key string, value any, expiration time.Duration) (bool, error) {
	return c.rdb.SetNX(ctx, key, value, expiration).Result()
}

func (c *client) Get(ctx context.Context, key string) (string, error) {
	return c.rdb.Get(ctx, key).Result()
}

func (c *client) Del(ctx context.Context, key string) error {
	return c.rdb.Del(ctx, key).Err()
}

func (c *client) Ping(ctx context.Context) error {
	return c.rdb.Ping(ctx).Err()
}

func (c *client) Close() error {
	return c.rdb.Close()
}

func (c *client) VerifyCaptcha(ctx context.Context, key, answer string) (string, error) {
	return verifyCaptchaScript.Run(ctx, c.rdb, []string{key}, answer).Text()
}
