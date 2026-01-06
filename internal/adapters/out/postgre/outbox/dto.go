package event_publisher

import (
	"time"

	"github.com/google/uuid"
)

type OutboxDTO struct {
	ID           uuid.UUID `db:"id"`
	EventName    string    `db:"event_name"`
	EventPayload string    `db:"event_payload"`
	OccurredAt   time.Time `db:"occurred_at"`
}
