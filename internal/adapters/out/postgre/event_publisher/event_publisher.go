package event_publisher

import (
	"context"
	"delivery/internal/pkg/ddd"
	"delivery/internal/pkg/outbox"

	"github.com/Masterminds/squirrel"
	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
)

const (
	tableNameOutbox    = "outbox"
	columnId           = "id"
	columnEventName    = "event_name"
	columnEventPayload = "event_payload"
	columnOccurredAt   = "occurred_at"
)

type txGetter interface {
	DefaultTrOrDB(ctx context.Context, db trmsqlx.Tr) trmsqlx.Tr
}

type eventEncoder interface {
	EncodeDomainEvent(event ddd.DomainEvent) (outbox.Message, error)
}

type EventPublisher struct {
	txGetter     txGetter
	eventEncoder eventEncoder
}

func NewEventPublisher(txGetter txGetter, eventEncoder eventEncoder) *EventPublisher {
	return &EventPublisher{txGetter: txGetter, eventEncoder: eventEncoder}
}

func (e *EventPublisher) Publish(ctx context.Context, event ddd.DomainEvent) error {
	tx := e.txGetter.DefaultTrOrDB(ctx, nil)

	outboxMessage, err := e.eventEncoder.EncodeDomainEvent(event)
	if err != nil {
		return err
	}

	outboxMessageDTO := convertOutboxMessageToDTO(outboxMessage)

	query, args, err := squirrel.Insert(tableNameOutbox).
		Columns(columnId, columnEventName, columnEventPayload, columnOccurredAt).
		Values(outboxMessageDTO.ID, outboxMessageDTO.EventName, outboxMessageDTO.EventPayload, outboxMessageDTO.OccurredAt).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func convertOutboxMessageToDTO(outboxMessage outbox.Message) OutboxDTO {
	return OutboxDTO{
		ID:           outboxMessage.ID,
		EventName:    outboxMessage.Name,
		EventPayload: string(outboxMessage.Payload),
		OccurredAt:   outboxMessage.OccurredAtUtc,
	}
}
