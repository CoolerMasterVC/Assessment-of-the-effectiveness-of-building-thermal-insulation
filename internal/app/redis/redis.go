// internal/app/redis/redis.go
package redis

import (
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
	"golang.org/x/net/context"
)

type Client struct {
	client *redis.Client
}

func NewRedisClient(addr, password string, db int) *Client {
	return &Client{
		client: redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: password,
			DB:       db,
		}),
	}
}

func (r *Client) SetSession(ctx context.Context, sessionID string, userID uint, expiration time.Duration) error {
	sessionData := map[string]interface{}{
		"user_id": userID,
		"created": time.Now(),
	}

	data, err := json.Marshal(sessionData)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, "session:"+sessionID, data, expiration).Err()
}

func (r *Client) GetSession(ctx context.Context, sessionID string) (uint, error) {
	data, err := r.client.Get(ctx, "session:"+sessionID).Result()
	if err != nil {
		return 0, err
	}

	var sessionData map[string]interface{}
	if err := json.Unmarshal([]byte(data), &sessionData); err != nil {
		return 0, err
	}

	return uint(sessionData["user_id"].(float64)), nil
}

func (r *Client) DeleteSession(ctx context.Context, sessionID string) error {
	return r.client.Del(ctx, "session:"+sessionID).Err()
}
