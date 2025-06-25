package ddd

type BaseAggregate[ID comparable] struct {
	baseEntity   *BaseEntity[ID]
	domainEvents []DomainEvent
}

func NewBaseAggregate[ID comparable](id ID) *BaseAggregate[ID] {
	return &BaseAggregate[ID]{
		baseEntity:   NewBaseEntity[ID](id),
		domainEvents: make([]DomainEvent, 0),
	}
}

func (ba *BaseAggregate[ID]) ID() ID {
	return ba.baseEntity.ID()
}

func (ba *BaseAggregate[ID]) Equal(other *BaseAggregate[ID]) bool {
	if other == nil {
		return false
	}
	return ba.baseEntity.Equal(other.baseEntity)
}

func (ba *BaseAggregate[ID]) ClearDomainEvents() {
	ba.domainEvents = []DomainEvent{}
}

func (ba *BaseAggregate[ID]) GetDomainEvents() []DomainEvent {
	return ba.domainEvents
}

func (ba *BaseAggregate[ID]) RaiseDomainEvent(event DomainEvent) {
	ba.domainEvents = append(ba.domainEvents, event)
}
