package application

import (
	"nc-two/domain/entity"
	"nc-two/domain/repository"
)

type postApp struct {
	fr repository.PostRepository
}

var _ PostAppInterface = &postApp{}

type PostAppInterface interface {
	SavePost(*entity.Post) (*entity.Post, map[string]string)
	GetAllPost() ([]entity.Post, error)
	GetPost(uint64) (*entity.Post, error)
	UpdatePost(*entity.Post) (*entity.Post, map[string]string)
	DeletePost(uint64) error
}

func (f *postApp) SavePost(post *entity.Post) (*entity.Post, map[string]string) {
	return f.fr.SavePost(post)
}

func (f *postApp) GetAllPost() ([]entity.Post, error) {
	return f.fr.GetAllPost()
}

func (f *postApp) GetPost(postId uint64) (*entity.Post, error) {
	return f.fr.GetPost(postId)
}

func (f *postApp) UpdatePost(post *entity.Post) (*entity.Post, map[string]string) {
	return f.fr.UpdatePost(post)
}

func (f *postApp) DeletePost(postId uint64) error {
	return f.fr.DeletePost(postId)
}
