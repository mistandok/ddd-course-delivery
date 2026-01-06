package event_publisher

import (
	"context"
	"delivery/internal/pkg/ddd"
	"delivery/internal/pkg/outbox"
	"time"

	"github.com/Masterminds/squirrel"
	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

const (
	tableNameOutbox    = "outbox"
	columnId           = "id"
	columnEventName    = "event_name"
	columnEventPayload = "event_payload"
	columnOccurredAt   = "occurred_at"
	columnProcessedAt  = "processed_at"
)

type txGetter interface {
	DefaultTrOrDB(ctx context.Context, db trmsqlx.Tr) trmsqlx.Tr
}

type eventEncoder interface {
	EncodeDomainEvent(event ddd.DomainEvent) (outbox.Message, error)
}

type Repository struct {
	db           *sqlx.DB
	txGetter     txGetter
	eventEncoder eventEncoder
}

func NewRepository(db *sqlx.DB, txGetter txGetter, eventEncoder eventEncoder) *Repository {
	return &Repository{db: db, txGetter: txGetter, eventEncoder: eventEncoder}
}

func (e *Repository) Publish(ctx context.Context, event ddd.DomainEvent) error {
	tx := e.txGetter.DefaultTrOrDB(ctx, e.db)

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

func (e *Repository) ProcessMessage(ctx context.Context, messageID uuid.UUID) error {
	tx := e.txGetter.DefaultTrOrDB(ctx, e.db)

	query, args, err := squirrel.Update(tableNameOutbox).
		Where(squirrel.Eq{columnId: messageID}).
		Set(columnProcessedAt, time.Now().UTC()).
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

func (e *Repository) GetUnprocessedMessages(ctx context.Context, limit uint64) ([]outbox.Message, error) {
	tx := e.txGetter.DefaultTrOrDB(ctx, e.db)

	query, args, err := squirrel.Select(columnId, columnEventName, columnEventPayload, columnOccurredAt, columnProcessedAt).
		From(tableNameOutbox).
		Where(squirrel.Eq{columnProcessedAt: nil}).
		Limit(limit).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	messages := make([]outbox.Message, 0)
	for rows.Next() {
		var message outbox.Message
		err = rows.Scan(&message.ID, &message.Name, &message.Payload, &message.OccurredAtUtc, &message.ProcessedAtUtc)
		if err != nil {
		}
		messages = append(messages, message)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}

func convertOutboxMessageToDTO(outboxMessage outbox.Message) OutboxDTO {
	return OutboxDTO{
		ID:           outboxMessage.ID,
		EventName:    outboxMessage.Name,
		EventPayload: string(outboxMessage.Payload),
		OccurredAt:   outboxMessage.OccurredAtUtc,
	}
}
