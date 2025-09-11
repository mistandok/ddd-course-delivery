package ports

//go:generate mockery --name UnitOfWorkFactory --with-expecter --exported
type UnitOfWorkFactory interface {
	NewUOW() UnitOfWork
}
