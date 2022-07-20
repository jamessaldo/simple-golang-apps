package unit

import (
	"nc-two/adapters"
	"nc-two/domain/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

//IF YOU HAVE TIME, YOU CAN TEST ALL THE METHODS FAILURES

var (
	saveUserRepo                func(*models.User) (*models.User, map[string]string)
	getUserRepo                 func(userId uint64) (*models.User, error)
	getUsersRepo                func() ([]models.User, error)
	getUserEmailAndPasswordRepo func(*models.User) (*models.User, map[string]string)
)

type fakeUserRepo struct{}

func (u *fakeUserRepo) SaveUser(user *models.User) (*models.User, map[string]string) {
	return saveUserRepo(user)
}
func (u *fakeUserRepo) GetUser(userId uint64) (*models.User, error) {
	return getUserRepo(userId)
}
func (u *fakeUserRepo) GetUsers() ([]models.User, error) {
	return getUsersRepo()
}
func (u *fakeUserRepo) GetUserByEmailAndPassword(user *models.User) (*models.User, map[string]string) {
	return getUserEmailAndPasswordRepo(user)
}

var userRepoFake adapters.UserRepository = &fakeUserRepo{} //this is where the real implementation is swap with our fake implementation

func TestSaveUser_Success(t *testing.T) {
	//Mock the response coming from the infrastructure
	saveUserRepo = func(user *models.User) (*models.User, map[string]string) {
		return &models.User{
			ID:        1,
			FirstName: "victor",
			LastName:  "steven",
			Email:     "steven@example.com",
			Password:  "password",
		}, nil
	}
	user := &models.User{
		ID:        1,
		FirstName: "victor",
		LastName:  "steven",
		Email:     "steven@example.com",
		Password:  "password",
	}
	u, err := userRepoFake.SaveUser(user)
	assert.Nil(t, err)
	assert.EqualValues(t, u.FirstName, "victor")
	assert.EqualValues(t, u.LastName, "steven")
	assert.EqualValues(t, u.Email, "steven@example.com")
}

func TestGetUser_Success(t *testing.T) {
	//Mock the response coming from the infrastructure
	getUserRepo = func(userId uint64) (*models.User, error) {
		return &models.User{
			ID:        1,
			FirstName: "victor",
			LastName:  "steven",
			Email:     "steven@example.com",
			Password:  "password",
		}, nil
	}
	userId := uint64(1)
	u, err := userRepoFake.GetUser(userId)
	assert.Nil(t, err)
	assert.EqualValues(t, u.FirstName, "victor")
	assert.EqualValues(t, u.LastName, "steven")
	assert.EqualValues(t, u.Email, "steven@example.com")
}

func TestGetUsers_Success(t *testing.T) {
	//Mock the response coming from the infrastructure
	getUsersRepo = func() ([]models.User, error) {
		return []models.User{
			{
				ID:        1,
				FirstName: "victor",
				LastName:  "steven",
				Email:     "steven@example.com",
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
	getUserEmailAndPasswordRepo = func(user *models.User) (*models.User, map[string]string) {
		return &models.User{
			ID:        1,
			FirstName: "victor",
			LastName:  "steven",
			Email:     "steven@example.com",
			Password:  "password",
		}, nil
	}
	user := &models.User{
		ID:        1,
		FirstName: "victor",
		LastName:  "steven",
		Email:     "steven@example.com",
		Password:  "password",
	}
	u, err := userRepoFake.GetUserByEmailAndPassword(user)
	assert.Nil(t, err)
	assert.EqualValues(t, u.FirstName, "victor")
	assert.EqualValues(t, u.LastName, "steven")
	assert.EqualValues(t, u.Email, "steven@example.com")
}
