package ddd

type BaseAggregate[ID comparable] struct {
	*BaseEntity[ID]
	domainEvents []DomainEvent
}

func NewBaseAggregate[ID comparable](id ID) *BaseAggregate[ID] {
	return &BaseAggregate[ID]{
		BaseEntity:   NewBaseEntity[ID](id),
		domainEvents: make([]DomainEvent, 0),
	}
}

func (a *BaseAggregate[ID]) ClearDomainEvents() {
	a.domainEvents = []DomainEvent{}
}

func (a *BaseAggregate[ID]) GetDomainEvents() []DomainEvent {
	return a.domainEvents
}

func (a *BaseAggregate[ID]) RaiseDomainEvent(event DomainEvent) {
	a.domainEvents = append(a.domainEvents, event)
}
