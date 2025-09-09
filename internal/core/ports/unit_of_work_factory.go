package ports

type UnitOfWorkFactory interface {
	NewUOW() (UnitOfWork, error)
}
