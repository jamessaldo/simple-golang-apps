package application

import (
	"nctwo/backend/adapters"
	"nctwo/backend/domain"
)

type PostApp struct {
	fr adapters.PostRepository
}

var _ PostAppInterface = &PostApp{}

type PostAppInterface interface {
	SavePost(*domain.Post) (*domain.Post, map[string]string)
	GetAllPost() ([]domain.Post, error)
	GetPost(uint64) (*domain.Post, error)
	UpdatePost(*domain.Post) (*domain.Post, map[string]string)
	DeletePost(uint64) error
}

func (p *PostApp) SavePost(post *domain.Post) (*domain.Post, map[string]string) {
	return p.fr.SavePost(post)
}

func (p *PostApp) GetAllPost() ([]domain.Post, error) {
	return p.fr.GetAllPost()
}

func (p *PostApp) GetPost(postId uint64) (*domain.Post, error) {
	return p.fr.GetPost(postId)
}

func (p *PostApp) UpdatePost(post *domain.Post) (*domain.Post, map[string]string) {
	return p.fr.UpdatePost(post)
}

func (p *PostApp) DeletePost(postId uint64) error {
	return p.fr.DeletePost(postId)
}
