package application

import (
	"nc-two/adapters"
	"nc-two/domain"
)

type postApp struct {
	fr adapters.PostRepository
}

var _ PostAppInterface = &postApp{}

type PostAppInterface interface {
	SavePost(*domain.Post) (*domain.Post, map[string]string)
	GetAllPost() ([]domain.Post, error)
	GetPost(uint64) (*domain.Post, error)
	UpdatePost(*domain.Post) (*domain.Post, map[string]string)
	DeletePost(uint64) error
}

func (f *postApp) SavePost(post *domain.Post) (*domain.Post, map[string]string) {
	return f.fr.SavePost(post)
}

func (f *postApp) GetAllPost() ([]domain.Post, error) {
	return f.fr.GetAllPost()
}

func (f *postApp) GetPost(postId uint64) (*domain.Post, error) {
	return f.fr.GetPost(postId)
}

func (f *postApp) UpdatePost(post *domain.Post) (*domain.Post, map[string]string) {
	return f.fr.UpdatePost(post)
}

func (f *postApp) DeletePost(postId uint64) error {
	return f.fr.DeletePost(postId)
}
