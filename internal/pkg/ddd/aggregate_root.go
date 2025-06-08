package ddd

type AggregateRoot interface {
	GetDomainEvents() []DomainEvent
	ClearDomainEvents()
	RaiseDomainEvent(DomainEvent)
}
