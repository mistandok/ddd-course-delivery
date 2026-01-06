package crons

import (
	"context"
	"log"

	"delivery/internal/pkg/ddd"
	"delivery/internal/pkg/errs"
	"delivery/internal/pkg/outbox"

	"github.com/google/uuid"
	"github.com/robfig/cron/v3"
)

const (
	limitUnprocessedMessages = 100
)

var _ cron.Job = &OutboxMessagesJob{}

type eventPublisher interface {
	Publish(ctx context.Context, domainEvent ddd.DomainEvent) error
}

type outboxRepository interface {
	GetUnprocessedMessages(ctx context.Context, limit uint64) ([]outbox.Message, error)
	ProcessMessage(ctx context.Context, messageID uuid.UUID) error
}

type eventRegistry interface {
	DecodeDomainEvent(message *outbox.Message) (ddd.DomainEvent, error)
}

type OutboxMessagesJob struct {
	eventPublisher   eventPublisher
	outboxRepository outboxRepository
	eventRegistry    eventRegistry
}

func NewOutboxMessagesJob(
	eventPublisher eventPublisher,
	outboxRepository outboxRepository,
	eventRegistry eventRegistry) (cron.Job, error) {
	if eventPublisher == nil {
		return nil, errs.NewValueIsRequiredError("eventPublisher")
	}

	return &OutboxMessagesJob{
		eventPublisher:   eventPublisher,
		outboxRepository: outboxRepository,
		eventRegistry:    eventRegistry,
	}, nil
}

func (j *OutboxMessagesJob) Run() {
	ctx := context.Background()
	messages, err := j.outboxRepository.GetUnprocessedMessages(ctx, limitUnprocessedMessages)
	if err != nil {
		log.Printf("OutboxMessagesJob error: %v", err)
		return
	}

	for _, message := range messages {
		domainEvent, err := j.eventRegistry.DecodeDomainEvent(&message)
		if err != nil {
			log.Printf("OutboxMessagesJob error: %v", err)
			continue
		}

		err = j.eventPublisher.Publish(ctx, domainEvent)
		if err != nil {
			log.Printf("OutboxMessagesJob error: %v", err)
			continue
		}

		err = j.outboxRepository.ProcessMessage(ctx, message.ID)
		if err != nil {
			log.Printf("OutboxMessagesJob error: %v", err)
		}
	}
}
