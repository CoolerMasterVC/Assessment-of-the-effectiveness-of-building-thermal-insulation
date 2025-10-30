package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type Client struct {
	client *redis.Client
}

func New(host string, port int) (*Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return &Client{client: rdb}, nil
}

func (c *Client) SetJWTBlacklist(ctx context.Context, token string, expiration time.Duration) error {
	return c.client.Set(ctx, "blacklist:"+token, "1", expiration).Err()
}

func (c *Client) IsJWTBlacklisted(ctx context.Context, token string) (bool, error) {
	val, err := c.client.Get(ctx, "blacklist:"+token).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return val == "1", nil
}

func (c *Client) Close() error {
	return c.client.Close()
}
