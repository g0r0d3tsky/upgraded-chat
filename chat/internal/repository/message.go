package repository

import (
	"2024-spring-ab-go-hw-3-g0r0d3tsky/chat/internal/domain"
	"context"

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

//
//func (s *StorageMessage) CreateMessage(ctx context.Context, mes *domain.Message) error {
//	id := uuid.New()
//	mes.ID = id
//	t := time.Now()
//	mes.Time = t
//	_, err := s.db.Exec(ctx,
//		`INSERT INTO "messages" (id, user_nickname, content, time) VALUES($1, $2, $3, $4)`,
//		&mes.ID, &mes.UserNickname, &mes.Content, &mes.Time,
//	)
//	if err != nil {
//		return err
//	}
//	return nil
//}

func (s *StorageMessage) GetAmountMessage(ctx context.Context, amount int) ([]*domain.Message, error) {
	var messages []*domain.Message

	rows, err := s.db.Query(
		ctx,
		`SELECT *
					FROM (
   						 SELECT id, user_nickname, content, time
    						FROM messages
   						 ORDER BY time DESC
   						 LIMIT $1
				) AS subquery
			ORDER BY time ASC;`, amount)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		mes := &domain.Message{}
		if err := rows.Scan(&mes.ID, &mes.UserNickname, &mes.Content, &mes.Time); err != nil {
			return nil, err
		}

		messages = append(messages, mes)
	}

	return messages, nil
}
