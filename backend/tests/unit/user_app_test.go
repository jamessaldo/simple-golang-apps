package unit

import (
	"nctwo/backend/application"
	"nctwo/backend/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

//IF YOU HAVE TIME, YOU CAN TEST ALL THE METHODS FAILURES

var (
	saveUserRepo                func(*domain.User) (*domain.User, map[string]string)
	getUserRepo                 func(userId uint64) (*domain.User, error)
	getUsersRepo                func() ([]domain.User, error)
	getUserEmailAndPasswordRepo func(*domain.User) (*domain.User, map[string]string)
)

type fakeUserRepo struct{}

func (u *fakeUserRepo) SaveUser(user *domain.User) (*domain.User, map[string]string) {
	return saveUserRepo(user)
}
func (u *fakeUserRepo) GetUser(userId uint64) (*domain.User, error) {
	return getUserRepo(userId)
}
func (u *fakeUserRepo) GetUsers() ([]domain.User, error) {
	return getUsersRepo()
}
func (u *fakeUserRepo) GetUserByEmailAndPassword(user *domain.User) (*domain.User, map[string]string) {
	return getUserEmailAndPasswordRepo(user)
}

var userRepoFake application.UserAppInterface = &fakeUserRepo{} //this is where the real implementation is swap with our fake implementation

func TestSaveUser_Success(t *testing.T) {
	//Mock the response coming from the infrastructure
	saveUserRepo = func(user *domain.User) (*domain.User, map[string]string) {
		return &domain.User{
			ID:        1,
			FirstName: "james",
			LastName:  "saldo",
			Email:     "jamessaldo@example.com",
			Password:  "password",
		}, nil
	}
	user := &domain.User{
		ID:        1,
		FirstName: "james",
		LastName:  "saldo",
		Email:     "jamessaldo@example.com",
		Password:  "password",
	}
	u, err := userRepoFake.SaveUser(user)
	assert.Nil(t, err)
	assert.EqualValues(t, u.FirstName, "james")
	assert.EqualValues(t, u.LastName, "saldo")
	assert.EqualValues(t, u.Email, "jamessaldo@example.com")
}

func TestGetUser_Success(t *testing.T) {
	//Mock the response coming from the infrastructure
	getUserRepo = func(userId uint64) (*domain.User, error) {
		return &domain.User{
			ID:        1,
			FirstName: "james",
			LastName:  "saldo",
			Email:     "jamessaldo@example.com",
			Password:  "password",
		}, nil
	}
	userId := uint64(1)
	u, err := userRepoFake.GetUser(userId)
	assert.Nil(t, err)
	assert.EqualValues(t, u.FirstName, "james")
	assert.EqualValues(t, u.LastName, "saldo")
	assert.EqualValues(t, u.Email, "jamessaldo@example.com")
}

func TestGetUsers_Success(t *testing.T) {
	//Mock the response coming from the infrastructure
	getUsersRepo = func() ([]domain.User, error) {
		return []domain.User{
			{
				ID:        1,
				FirstName: "james",
				LastName:  "saldo",
				Email:     "jamessaldo@example.com",
				Password:  "password",
			},
			{
				ID:        2,
				FirstName: "kobe",
				LastName:  "bryant",
				Email:     "kobe@example.com",
				Password:  "password",
			},
		}, nil
	}
	users, err := userRepoFake.GetUsers()
	assert.Nil(t, err)
	assert.EqualValues(t, len(users), 2)
}

func TestGetUserByEmailAndPassword_Success(t *testing.T) {
	//Mock the response coming from the infrastructure
	getUserEmailAndPasswordRepo = func(user *domain.User) (*domain.User, map[string]string) {
		return &domain.User{
			ID:        1,
			FirstName: "james",
			LastName:  "saldo",
			Email:     "jamessaldo@example.com",
			Password:  "password",
		}, nil
	}
	user := &domain.User{
		ID:        1,
		FirstName: "james",
		LastName:  "saldo",
		Email:     "jamessaldo@example.com",
		Password:  "password",
	}
	u, err := userRepoFake.GetUserByEmailAndPassword(user)
	assert.Nil(t, err)
	assert.EqualValues(t, u.FirstName, "james")
	assert.EqualValues(t, u.LastName, "saldo")
	assert.EqualValues(t, u.Email, "jamessaldo@example.com")
}
