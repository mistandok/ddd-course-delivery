package ddd

type BaseEntity[ID comparable] struct {
	id ID
}

func NewBaseEntity[ID comparable](id ID) *BaseEntity[ID] {
	return &BaseEntity[ID]{
		id: id,
	}
}

func (be *BaseEntity[ID]) ID() ID {
	return be.id
}

func (be *BaseEntity[ID]) Equal(other *BaseEntity[ID]) bool {
	if other == nil {
		return false
	}
	return be.id == other.id
}
