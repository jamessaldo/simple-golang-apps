package persistence

import (
	"fmt"
	"nctwo/backend/adapters"
	"nctwo/backend/domain"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Repositories struct {
	User    adapters.UserRepository
	Post    adapters.PostRepository
	Comment adapters.CommentRepository
	db      *gorm.DB
}

func NewRepositories(Dbdriver, DbUser, DbPassword, DbPort, DbHost, DbName string) (*Repositories, error) {
	DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", DbHost, DbPort, DbUser, DbName, DbPassword)
	db, err := gorm.Open(Dbdriver, DBURL)
	if err != nil {
		return nil, err
	}
	db.LogMode(true)

	return &Repositories{
		User:    adapters.NewUserRepository(db),
		Post:    adapters.NewPostRepository(db),
		Comment: adapters.NewCommentRepository(db),
		db:      db,
	}, nil
}

//closes the  database connection
func (s *Repositories) Close() error {
	return s.db.Close()
}

//This migrate all tables
func (s *Repositories) Automigrate() error {
	return s.db.AutoMigrate(&domain.User{}, &domain.Post{}, &domain.Comment{}).Error
}
