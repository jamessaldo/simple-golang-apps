package handlers

import (
	"nc-two/adapters"
	"nc-two/domain/models"
	"nc-two/infrastructure/auth"
)

//Users struct defines the dependencies that will be used
type Users struct {
	us adapters.UserRepository
	rd auth.AuthInterface
	tk auth.TokenInterface
}

//Users constructor
func NewUsers(us adapters.UserRepository, rd auth.AuthInterface, tk auth.TokenInterface) *Users {
	return &Users{
		us: us,
		rd: rd,
		tk: tk,
	}
}

func (s *Users) SaveUser(user *models.User) (interface{}, map[string]string) {
	newUser, err := s.us.SaveUser(user)
	if err != nil {
		return nil, err
	}
	return newUser.PublicUser(), nil
}

func (s *Users) GetUsers() (models.Users, error) {
	users := models.Users{} //customize user
	var err error
	users, err = s.us.GetUsers()
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *Users) GetUser(userId uint64) (*models.PublicUser, error) {
	user, err := s.us.GetUser(userId)
	if err != nil {
		return nil, err
	}
	return user.PublicUser(), nil
}
