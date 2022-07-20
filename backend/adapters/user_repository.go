package adapters

import (
	"errors"
	"nc-two/domain/models"
	"nc-two/infrastructure/security"
	"strings"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type UserRepo struct {
	db *gorm.DB
}

//UserRepo implements the UserRepository interface
type UserRepository interface {
	SaveUser(*models.User) (*models.User, map[string]string)
	GetUser(uint64) (*models.User, error)
	GetUsers() ([]models.User, error)
	GetUserByEmailAndPassword(*models.User) (*models.User, map[string]string)
}

func NewUserRepository(db *gorm.DB) *UserRepo {
	return &UserRepo{db}
}

var _ UserRepository = &UserRepo{}

func (r *UserRepo) SaveUser(user *models.User) (*models.User, map[string]string) {
	dbErr := map[string]string{}
	err := r.db.Debug().Create(&user).Error
	if err != nil {
		//If the email is already taken
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "Duplicate") {
			dbErr["email_taken"] = "email already taken"
			return nil, dbErr
		}
		//any other db error
		dbErr["db_error"] = "database error"
		return nil, dbErr
	}
	return user, nil
}

func (r *UserRepo) GetUser(id uint64) (*models.User, error) {
	var user models.User
	err := r.db.Debug().Where("id = ?", id).Take(&user).Error
	if err != nil {
		return nil, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return nil, errors.New("user not found")
	}
	return &user, nil
}

func (r *UserRepo) GetUsers() ([]models.User, error) {
	var users []models.User
	err := r.db.Debug().Find(&users).Error
	if err != nil {
		return nil, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return nil, errors.New("user not found")
	}
	return users, nil
}

func (r *UserRepo) GetUserByEmailAndPassword(u *models.User) (*models.User, map[string]string) {
	var user models.User
	dbErr := map[string]string{}
	err := r.db.Debug().Where("email = ?", u.Email).Take(&user).Error
	if gorm.IsRecordNotFoundError(err) {
		dbErr["no_user"] = "user not found"
		return nil, dbErr
	}
	if err != nil {
		dbErr["db_error"] = "database error"
		return nil, dbErr
	}
	//Verify the password
	err = security.VerifyPassword(user.Password, u.Password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		dbErr["incorrect_password"] = "incorrect password"
		return nil, dbErr
	}
	return &user, nil
}
