package repository

import (
	"context"
	"fmt"
	"time"

	"2024-spring-ab-go-hw-3-g0r0d3tsky/storage/internal/config"
	"2024-spring-ab-go-hw-3-g0r0d3tsky/storage/internal/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type StorageMessage struct {
	db *pgxpool.Pool
}

func NewStorageMessage(dbPool *pgxpool.Pool) StorageMessage {
	Storage := StorageMessage{
		db: dbPool,
	}
	return Storage
}

func (s *StorageMessage) CreateMessage(ctx context.Context, mes *domain.Message) error {
	id := uuid.New()
	mes.ID = id
	_, err := s.db.Exec(ctx,
		`INSERT INTO "messages" (id, user_nickname, content, time) VALUES($1, $2, $3, $4)`,
		&mes.ID, &mes.UserNickname, &mes.Content, &mes.Time,
	)
	if err != nil {
		return err
	}
	return nil
}

func Connect(c *config.Config) (*pgxpool.Pool, error) {
	connectionString := c.PostgresDSN()

	poolConfig, err := pgxpool.ParseConfig(connectionString)
	if err != nil {
		return nil, fmt.Errorf("parsing pgx pool config: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("connecting to postgres: %w", err)
	}

	err = pool.Ping(context.Background())
	if err != nil {
		return nil, fmt.Errorf("connecting to postgres: %w", err)
	}

	return pool, nil
}
