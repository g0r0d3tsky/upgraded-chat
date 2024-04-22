package domain

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID           uuid.UUID
	UserNickname string
	Content      string
	Time         time.Time
}
