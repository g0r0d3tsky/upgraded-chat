package usecase

import (
	"context"
	"fmt"
	"time"

	"2024-spring-ab-go-hw-3-g0r0d3tsky/chat/internal/cache"
	"2024-spring-ab-go-hw-3-g0r0d3tsky/chat/internal/domain"
	"2024-spring-ab-go-hw-3-g0r0d3tsky/chat/internal/kafka"
)

//go:generate mockgen -source=message.go -destination=mocks/BrokerCacheMock.go
type Broker interface {
	Push(topic string, message *domain.Message) error
}

type Cache interface {
	GetMessages(ctx context.Context, amount int) ([]*domain.Message, error)
	SetMessages(ctx context.Context, messages []*domain.Message) error
}

type MessageService struct {
	broker Broker
	cache  Cache
}

func NewMessageService(producer *kafka.Producer, cache cache.Cache) *MessageService {
	return &MessageService{
		broker: producer,
		cache:  cache,
	}
}

func (s *MessageService) Push(ctx context.Context, topic string, mess *domain.Message) error {
	t := time.Now()
	mess.Time = t
	err := s.broker.Push(topic, mess)
	if err != nil {
		return fmt.Errorf("pushing kafka: %w", err)
	}
	var messages []*domain.Message
	messages = append(messages, mess)
	err = s.cache.SetMessages(ctx, messages)
	if err != nil {
		return fmt.Errorf("setting message: %w", err)
	}
	return nil
}

func (s *MessageService) GetMessages(ctx context.Context, amount int) ([]*domain.Message, error) {
	return s.cache.GetMessages(ctx, amount)
}

func (s *MessageService) SetMessages(ctx context.Context, messages []*domain.Message) error {
	return s.cache.SetMessages(ctx, messages)
}
