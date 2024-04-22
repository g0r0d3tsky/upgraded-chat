package usecase

import (
	"2024-spring-ab-go-hw-3-g0r0d3tsky/storage/internal/domain"
	"context"
	"fmt"
)

type Repository interface {
	CreateMessage(ctx context.Context, mes *domain.Message) error
}

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) CreateMessage(ctx context.Context, message *domain.Message) error {
	err := s.repository.CreateMessage(ctx, message)
	if err != nil {
		return fmt.Errorf("creating message: %w", err)
	}
	return nil
}
