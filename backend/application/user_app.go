package application

import (
	"nctwo/backend/adapters"
	"nctwo/backend/domain"
)

type userApp struct {
	us adapters.UserRepository
}

//UserApp implements the UserAppInterface
var _ UserAppInterface = &userApp{}

type UserAppInterface interface {
	SaveUser(*domain.User) (*domain.User, map[string]string)
	GetUsers() ([]domain.User, error)
	GetUser(uint64) (*domain.User, error)
	GetUserByEmailAndPassword(*domain.User) (*domain.User, map[string]string)
}

func (u *userApp) SaveUser(user *domain.User) (*domain.User, map[string]string) {
	return u.us.SaveUser(user)
}

func (u *userApp) GetUser(userId uint64) (*domain.User, error) {
	return u.us.GetUser(userId)
}

func (u *userApp) GetUsers() ([]domain.User, error) {
	return u.us.GetUsers()
}

func (u *userApp) GetUserByEmailAndPassword(user *domain.User) (*domain.User, map[string]string) {
	return u.us.GetUserByEmailAndPassword(user)
}
