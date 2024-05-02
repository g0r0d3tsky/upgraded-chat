package cache

import (
	"context"
	"encoding/json"
	"fmt"

	"2024-spring-ab-go-hw-3-g0r0d3tsky/chat/internal/domain"

	"github.com/redis/go-redis/v9"
)

//go:generate mockgen -source=message.go -destination=mocks/repositoryMock.go
type Repository interface {
	GetAmountMessage(ctx context.Context, amount int) ([]*domain.Message, error)
}

type Redis interface {
	LLen(ctx context.Context, key string) *redis.IntCmd
	LRange(ctx context.Context, key string, start, stop int64) *redis.StringSliceCmd
	RPush(ctx context.Context, key string, values ...interface{}) *redis.IntCmd
}

type Cache struct {
	redis Redis
	repo  Repository
}

func ToRedisList(p []*domain.Message) []*Message {
	r := make([]*Message, len(p))

	for i := range p {
		r[i] = &Message{
			ID:       p[i].ID,
			Nickname: p[i].UserNickname,
			Content:  p[i].Content,
			Time:     p[i].Time,
		}
	}

	return r
}

func ToDomainList(p []*Message) []*domain.Message {
	r := make([]*domain.Message, len(p))

	for i := range p {
		r[i] = &domain.Message{
			ID:           p[i].ID,
			UserNickname: p[i].Nickname,
			Content:      p[i].Content,
			Time:         p[i].Time,
		}
	}
	return r
}

func NewCache(client Redis, repo Repository) Cache {
	cache := Cache{
		redis: client,
		repo:  repo,
	}
	return cache
}

func (c Cache) GetMessages(ctx context.Context, amount int) ([]*domain.Message, error) {
	length, err := c.redis.LLen(ctx, "messages").Result()
	if err != nil {
		return nil, err
	}

	var redisMess []*Message

	start := length - int64(amount)
	if start < 0 {
		messages, err := c.repo.GetAmountMessage(ctx, amount)
		if err != nil {
			return nil, fmt.Errorf("getting messages from db: %w", err)
		}
		err = c.SetMessages(ctx, messages)
		if err != nil {
			return nil, fmt.Errorf("setting messages: %w", err)
		}
		return messages, nil
	}

	mess, err := c.redis.LRange(ctx, "messages", start, length-1).Result()
	if err != nil {
		return nil, err
	}

	for _, val := range mess {
		var data *Message
		err = json.Unmarshal([]byte(val), &data)
		if err != nil {
			return nil, fmt.Errorf("unmarshal from json redis: %w", err)
		}
		redisMess = append(redisMess, data)
	}

	return ToDomainList(redisMess), nil
}

func (c Cache) SetMessages(ctx context.Context, messages []*domain.Message) error {
	mess := ToRedisList(messages)
	for _, val := range mess {
		data, err := json.Marshal(val)
		if err != nil {
			return fmt.Errorf("marshal json redis: %w", err)
		}
		err = c.redis.RPush(ctx, "messages", data).Err()
		if err != nil {
			return fmt.Errorf("pushing redis: %w", err)
		}
	}

	return nil
}
