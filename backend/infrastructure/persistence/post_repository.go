package persistence

import (
	"errors"
	"nc-two/domain/entity"
	"nc-two/domain/repository"
	"os"
	"strings"

	"github.com/jinzhu/gorm"
)

type PostRepo struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) *PostRepo {
	return &PostRepo{db}
}

//PostRepo implements the repository.PostRepository interface
var _ repository.PostRepository = &PostRepo{}

func (r *PostRepo) SavePost(post *entity.Post) (*entity.Post, map[string]string) {
	dbErr := map[string]string{}
	//The images are uploaded to digital ocean spaces. So we need to prepend the url. This might not be your use case, if you are not uploading image to Digital Ocean.
	post.PostImage = os.Getenv("DO_SPACES_URL") + post.PostImage

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

func (r *PostRepo) GetPost(id uint64) (*entity.Post, error) {
	var post entity.Post
	err := r.db.Debug().Where("id = ?", id).Take(&post).Error
	if err != nil {
		return nil, errors.New("database error, please try again")
	}
	if gorm.IsRecordNotFoundError(err) {
		return nil, errors.New("post not found")
	}
	return &post, nil
}

func (r *PostRepo) GetAllPost() ([]entity.Post, error) {
	var posts []entity.Post
	err := r.db.Debug().Limit(100).Order("created_at desc").Find(&posts).Error
	if err != nil {
		return nil, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return nil, errors.New("user not found")
	}
	return posts, nil
}

func (r *PostRepo) UpdatePost(post *entity.Post) (*entity.Post, map[string]string) {
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
	var post entity.Post
	err := r.db.Debug().Where("id = ?", id).Delete(&post).Error
	if err != nil {
		return errors.New("database error, please try again")
	}
	return nil
}
