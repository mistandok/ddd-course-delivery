package common

import (
	"context"
	"delivery/internal/core/ports"
	"delivery/internal/pkg/ddd"
	"delivery/internal/pkg/errs"
	"encoding/json"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/gogo/protobuf/proto"
)

// Интеграционное событие, которое будет отправлено в Kafka
type IntegrationEvent[TIntegrationEvent proto.Message] struct {
	event TIntegrationEvent
	key   string
}

func NewIntegrationEvent[TIntegrationEvent proto.Message](event TIntegrationEvent, key string) *IntegrationEvent[TIntegrationEvent] {
	return &IntegrationEvent[TIntegrationEvent]{
		event: event,
		key:   key,
	}
}

func (e *IntegrationEvent[TIntegrationEvent]) Event() TIntegrationEvent {
	return e.event
}

func (e *IntegrationEvent[TIntegrationEvent]) Key() string {
	return e.key
}

// Маппер для преборазования доменного события в интеграционное событие
type FromDomainToIntegrationMapper[TDomainEvent ddd.DomainEvent, TIntegrationEvent proto.Message] interface {
	Map(domainEvent TDomainEvent) IntegrationEvent[TIntegrationEvent]
}

// Продюсер для отправки событий в Kafka
type KafkaProducer[TDomainEvent ddd.DomainEvent, TIntegrationEvent proto.Message] struct {
	topic       string
	producer    sarama.SyncProducer
	eventMapper FromDomainToIntegrationMapper[TDomainEvent, TIntegrationEvent]
}

func NewKafkaProducer[
	TDomainEvent ddd.DomainEvent,
	TIntegrationEvent proto.Message,
	TEventMapper FromDomainToIntegrationMapper[TDomainEvent, TIntegrationEvent],
](brokers []string, topic string, eventMapper TEventMapper) (ports.EventProducer[TDomainEvent], error) {
	if len(brokers) == 0 {
		return nil, errs.NewValueIsRequiredError("brokers")
	}
	if topic == "" {
		return nil, errs.NewValueIsRequiredError("topic")
	}

	saramaCfg := sarama.NewConfig()
	saramaCfg.Version = sarama.V3_4_0_0
	saramaCfg.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokers, saramaCfg)
	if err != nil {
		return nil, fmt.Errorf("create sync producer: %w", err)
	}

	return &KafkaProducer[TDomainEvent, TIntegrationEvent]{
		topic:       topic,
		eventMapper: eventMapper,
		producer:    producer,
	}, nil
}

func (p *KafkaProducer[TDomainEvent, TIntegrationEvent]) Close() error {
	return p.producer.Close()
}

func (p *KafkaProducer[TDomainEvent, TIntegrationEvent]) Publish(ctx context.Context, domainEvent TDomainEvent) error {
	integrationEvent := p.eventMapper.Map(domainEvent)

	bytes, err := json.Marshal(integrationEvent.Event())
	if err != nil {
		return fmt.Errorf("marshal event: %w", err)
	}

	msg := &sarama.ProducerMessage{
		Topic: p.topic,
		Key:   sarama.StringEncoder(integrationEvent.Key()),
		Value: sarama.ByteEncoder(bytes),
	}

	resultCh := make(chan error, 1)

	go func() {
		_, _, err := p.producer.SendMessage(msg)
		resultCh <- err
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-resultCh:
		return err
	}
}
