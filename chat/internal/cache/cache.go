package cache

import (
	"2024-spring-ab-go-hw-3-g0r0d3tsky/chat/internal/config"
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type Message struct {
	ID       uuid.UUID `json:"id"`
	Nickname string    `json:"nickname"`
	Content  string    `json:"content"`
	Time     time.Time `json:"time"`
}

func Connect(c *config.Config) (*redis.Client, error) {
	connectionString := c.RedisDSN()
	opts, err := redis.ParseURL(connectionString)
	if err != nil {
		slog.Error("parsing redis url")
	}

	client := redis.NewClient(opts)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	status, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("redis connection was refused %w", err)
	}
	if status != "PONG" {
		return nil, fmt.Errorf("redis connection was not successful")
	}

	return client, nil
}
