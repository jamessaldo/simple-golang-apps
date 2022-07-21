package application

import (
	"nc-two/adapters"
	"nc-two/domain"
)

type CommentApp struct {
	cr adapters.CommentRepository
}

var _ CommentAppInterface = &CommentApp{}

type CommentAppInterface interface {
	SaveComment(*domain.Comment) (*domain.Comment, map[string]string)
	GetAllComment() ([]domain.Comment, error)
	GetComment(uint64) (*domain.Comment, error)
	UpdateComment(*domain.Comment) (*domain.Comment, map[string]string)
	DeleteComment(uint64) error
}

func (c *CommentApp) SaveComment(comment *domain.Comment) (*domain.Comment, map[string]string) {
	return c.cr.SaveComment(comment)
}

func (c *CommentApp) GetAllComment() ([]domain.Comment, error) {
	return c.cr.GetAllComment()
}

func (c *CommentApp) GetComment(commentId uint64) (*domain.Comment, error) {
	return c.cr.GetComment(commentId)
}

func (c *CommentApp) UpdateComment(comment *domain.Comment) (*domain.Comment, map[string]string) {
	return c.cr.UpdateComment(comment)
}

func (c *CommentApp) DeleteComment(commentId uint64) error {
	return c.cr.DeleteComment(commentId)
}
