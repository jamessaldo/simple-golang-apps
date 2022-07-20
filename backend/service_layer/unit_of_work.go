package service_layer

import (
	"nc-two/adapters"
	"nc-two/infrastructure/persistence"
)

//UnitOfWork struct defines the dependencies that will be used
type UnitOfWork struct {
	Users adapters.UserRepository
	Posts adapters.PostRepository
}

//UnitOfWork constructor
func CreateUnitOfWork(s persistence.Repositories) *UnitOfWork {
	return &UnitOfWork{
		Users: s.User,
		Posts: s.Post,
	}
}
