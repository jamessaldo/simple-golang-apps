package adapters

import (
	"errors"
	"nc-two/domain"
	"strings"

	"github.com/jinzhu/gorm"
)

type PostRepo struct {
	db *gorm.DB
}

//PostRepo implements the PostRepository interface
type PostRepository interface {
	SavePost(*domain.Post) (*domain.Post, map[string]string)
	GetPost(uint64) (*domain.Post, error)
	GetAllPost() ([]domain.Post, error)
	UpdatePost(*domain.Post) (*domain.Post, map[string]string)
	DeletePost(uint64) error
}

func NewPostRepository(db *gorm.DB) *PostRepo {
	return &PostRepo{db}
}

var _ PostRepository = &PostRepo{}

func (r *PostRepo) SavePost(post *domain.Post) (*domain.Post, map[string]string) {
	dbErr := map[string]string{}

	err := r.db.Debug().Create(&post).Error
	if err != nil {
		//since our title is unique
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "Duplicate") {
			dbErr["unique_title"] = "post title already taken"
			return nil, dbErr
		}
		//any other db error
		dbErr["db_error"] = "database error"
		return nil, dbErr
	}
	return post, nil
}

func (r *PostRepo) GetPost(id uint64) (*domain.Post, error) {
	var post domain.Post
	err := r.db.Debug().Where("id = ?", id).Take(&post).Error
	if err != nil {
		return nil, errors.New("database error, please try again")
	}
	if gorm.IsRecordNotFoundError(err) {
		return nil, errors.New("post not found")
	}
	return &post, nil
}

func (r *PostRepo) GetAllPost() ([]domain.Post, error) {
	var posts []domain.Post
	err := r.db.Debug().Limit(100).Order("created_at desc").Find(&posts).Error
	if err != nil {
		return nil, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return nil, errors.New("user not found")
	}
	return posts, nil
}

func (r *PostRepo) UpdatePost(post *domain.Post) (*domain.Post, map[string]string) {
	dbErr := map[string]string{}
	err := r.db.Debug().Save(&post).Error
	if err != nil {
		//since our title is unique
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "Duplicate") {
			dbErr["unique_title"] = "title already taken"
			return nil, dbErr
		}
		//any other db error
		dbErr["db_error"] = "database error"
		return nil, dbErr
	}
	return post, nil
}

func (r *PostRepo) DeletePost(id uint64) error {
	var post domain.Post
	err := r.db.Debug().Where("id = ?", id).Delete(&post).Error
	if err != nil {
		return errors.New("database error, please try again")
	}
	return nil
}
