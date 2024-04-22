package handlers

import (
	"2024-spring-ab-go-hw-3-g0r0d3tsky/storage/internal/domain"
	"2024-spring-ab-go-hw-3-g0r0d3tsky/storage/internal/kafka"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/IBM/sarama"
)

type Service interface {
	CreateMessage(ctx context.Context, message *domain.Message) error
}

type ServiceMessage struct {
	service Service
}

func NewServiceMessage(service Service) ServiceMessage {
	return ServiceMessage{
		service: service,
	}
}

func (s *ServiceMessage) Handler() *kafka.Consumer {
	handler := func(msg *sarama.ConsumerMessage) error {
		var message *domain.Message
		err := json.Unmarshal(msg.Value, &message)
		if err != nil {
			slog.Error("unmarshalling json: %v", err)
			return fmt.Errorf("unmarshalling json: %w", err)
		}

		err = s.CreateMessage(context.Background(), message)
		if err != nil {
			slog.Error("creating message: %v", err)
			return fmt.Errorf("creating message: %w", err)
		}

		return nil
	}

	return kafka.NewConsumer(handler)
}

func (s *ServiceMessage) CreateMessage(ctx context.Context, message *domain.Message) error {
	err := s.service.CreateMessage(ctx, message)
	if err != nil {
		return fmt.Errorf("creating message: %w", err)
	}
	return nil
}
