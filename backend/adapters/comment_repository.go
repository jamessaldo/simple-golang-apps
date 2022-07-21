package adapters

import (
	"errors"
	"nc-two/domain"
	"strings"

	"github.com/jinzhu/gorm"
)

type CommentRepo struct {
	db *gorm.DB
}

//CommentRepo implements the CommentRepository interface
type CommentRepository interface {
	SaveComment(*domain.Comment) (*domain.Comment, map[string]string)
	GetComment(uint64) (*domain.Comment, error)
	GetAllComment() ([]domain.Comment, error)
	UpdateComment(*domain.Comment) (*domain.Comment, map[string]string)
	DeleteComment(uint64) error
}

func NewCommentRepository(db *gorm.DB) *CommentRepo {
	return &CommentRepo{db}
}

var _ CommentRepository = &CommentRepo{}

func (r *CommentRepo) SaveComment(comment *domain.Comment) (*domain.Comment, map[string]string) {
	dbErr := map[string]string{}

	err := r.db.Debug().Create(&comment).Error
	if err != nil {
		//any other db error
		dbErr["db_error"] = "database error"
		return nil, dbErr
	}
	return comment, nil
}

func (r *CommentRepo) GetComment(id uint64) (*domain.Comment, error) {
	var comment domain.Comment
	err := r.db.Debug().Where("id = ?", id).Take(&comment).Error
	if err != nil {
		return nil, errors.New("database error, please try again")
	}
	if gorm.IsRecordNotFoundError(err) {
		return nil, errors.New("comment not found")
	}
	return &comment, nil
}

func (r *CommentRepo) GetAllComment() ([]domain.Comment, error) {
	var comments []domain.Comment
	err := r.db.Debug().Limit(100).Order("created_at desc").Find(&comments).Error
	if err != nil {
		return nil, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return nil, errors.New("user not found")
	}
	return comments, nil
}

func (r *CommentRepo) UpdateComment(comment *domain.Comment) (*domain.Comment, map[string]string) {
	dbErr := map[string]string{}
	err := r.db.Debug().Save(&comment).Error
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
	return comment, nil
}

func (r *CommentRepo) DeleteComment(id uint64) error {
	var comment domain.Comment
	err := r.db.Debug().Where("id = ?", id).Delete(&comment).Error
	if err != nil {
		return errors.New("database error, please try again")
	}
	return nil
}
