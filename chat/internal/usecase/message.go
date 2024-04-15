package usecase

import (
	"2024-spring-ab-go-hw-3-g0r0d3tsky/chat/internal/domain"
	"context"
)

type MessageRepo interface {
	CreateMessage(ctx context.Context, mes *domain.Message) error
	GetAmountMessage(ctx context.Context, amount int) ([]*domain.Message, error)
}

type MessageService struct {
	repo MessageRepo
}

func NewMessageService(repo MessageRepo) *MessageService {
	return &MessageService{repo: repo}
}

func (s *MessageService) CreateMessage(ctx context.Context, message *domain.Message) error {
	return s.repo.CreateMessage(ctx, message)
}

func (s *MessageService) GetAmountMessage(ctx context.Context, amount int) ([]*domain.Message, error) {
	return s.repo.GetAmountMessage(ctx, amount)
}
