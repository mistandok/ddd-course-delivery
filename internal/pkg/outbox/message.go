package outbox

import (
	"github.com/google/uuid"
	"time"
)

type Message struct {
	ID             uuid.UUID
	Name           string
	Payload        []byte
	OccurredAtUtc  time.Time
	ProcessedAtUtc *time.Time
}

func (Message) TableName() string {
	return "outbox"
}
