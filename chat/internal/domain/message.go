package domain

import (
	"github.com/google/uuid"
	"time"
)

type Message struct {
	ID           uuid.UUID
	UserNickname string
	Content      string
	Time         time.Time
}
