package repository

import "nc-two/domain/entity"

type PostRepository interface {
	SavePost(*entity.Post) (*entity.Post, map[string]string)
	GetPost(uint64) (*entity.Post, error)
	GetAllPost() ([]entity.Post, error)
	UpdatePost(*entity.Post) (*entity.Post, map[string]string)
	DeletePost(uint64) error
}
