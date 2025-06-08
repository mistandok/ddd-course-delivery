package ddd

type BaseEntity[ID comparable] struct {
	id ID
}

func NewBaseEntity[ID comparable](id ID) *BaseEntity[ID] {
	return &BaseEntity[ID]{
		id: id,
	}
}

func (a *BaseEntity[ID]) ID() ID {
	return a.id
}

func (a *BaseEntity[ID]) Equal(other *BaseEntity[ID]) bool {
	if other == nil {
		return false
	}
	return a.id == other.id
}
