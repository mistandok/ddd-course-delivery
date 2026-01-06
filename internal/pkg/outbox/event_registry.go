package outbox

import (
	"fmt"
	"reflect"
	"time"

	"delivery/internal/pkg/ddd"
	"delivery/internal/pkg/errs"
	"delivery/internal/pkg/outbox/mappers"
)

type EventRegistry interface {
	RegisterDomainEvent(eventType ddd.DomainEvent, toMapper mappers.ToJSONMapper, fromMapper mappers.FromJSONMapper) error
	EncodeDomainEvent(domainEvent ddd.DomainEvent) (Message, error)
	DecodeDomainEvent(event *Message) (ddd.DomainEvent, error)
}

var _ EventRegistry = &eventRegistry{}

type eventRegistryEntry struct {
	toJSONMapper   mappers.ToJSONMapper
	fromJSONMapper mappers.FromJSONMapper
}

func (e *eventRegistryEntry) ToJSON(event ddd.DomainEvent) ([]byte, error) {
	return e.toJSONMapper.ToJSON(event)
}

func (e *eventRegistryEntry) FromJSON(data []byte) (ddd.DomainEvent, error) {
	return e.fromJSONMapper.FromJSON(data)
}

type eventRegistry struct {
	eventRegistry map[string]*eventRegistryEntry
}

func NewEventRegistry() (EventRegistry, error) {
	return &eventRegistry{
		eventRegistry: make(map[string]*eventRegistryEntry),
	}, nil
}

func (r *eventRegistry) RegisterDomainEvent(eventType ddd.DomainEvent, toMapper mappers.ToJSONMapper, fromMapper mappers.FromJSONMapper) error {
	if eventType == nil {
		return errs.NewValueIsRequiredError("eventType")
	}
	if toMapper == nil {
		return errs.NewValueIsRequiredError("toMapper")
	}
	if fromMapper == nil {
		return errs.NewValueIsRequiredError("fromMapper")
	}

	eventTypeName := r.getEventTypeName(eventType)

	r.eventRegistry[eventTypeName] = &eventRegistryEntry{
		toJSONMapper:   toMapper,
		fromJSONMapper: fromMapper,
	}
	return nil
}

func (r *eventRegistry) DecodeDomainEvent(outboxMessage *Message) (ddd.DomainEvent, error) {
	entry, ok := r.eventRegistry[outboxMessage.Name]
	if !ok {
		return nil, fmt.Errorf("unknown outboxMessage type: %s", outboxMessage.Name)
	}

	return entry.FromJSON(outboxMessage.Payload)
}

func (r *eventRegistry) EncodeDomainEvent(domainEvent ddd.DomainEvent) (Message, error) {
	eventTypeName := r.getEventTypeName(domainEvent)

	entry, ok := r.eventRegistry[eventTypeName]
	if !ok {
		return Message{}, fmt.Errorf("unknown domain event type: %s", eventTypeName)
	}

	payload, err := entry.ToJSON(domainEvent)
	if err != nil {
		return Message{}, fmt.Errorf("failed to marshal event: %w", err)
	}

	return Message{
		ID:             domainEvent.GetID(),
		Name:           eventTypeName,
		Payload:        payload,
		OccurredAtUtc:  time.Now().UTC(),
		ProcessedAtUtc: nil,
	}, nil
}

func (r *eventRegistry) getEventTypeName(eventType ddd.DomainEvent) string {
	t := reflect.TypeOf(eventType)

	// Если передан указатель, получаем тип элемента
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	return t.Name()
}
