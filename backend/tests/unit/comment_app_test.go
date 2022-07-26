package unit

import (
	"nctwo/backend/application"
	"nctwo/backend/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

//IF YOU HAVE TIME, YOU CAN TEST ALL THE METHODS FAILURES

type fakeCommentRepo struct{}

var (
	saveCommentRepo   func(*domain.Comment) (*domain.Comment, map[string]string)
	getCommentRepo    func(uint64) (*domain.Comment, error)
	getAllCommentRepo func() ([]domain.Comment, error)
	updateCommentRepo func(*domain.Comment) (*domain.Comment, map[string]string)
	deleteCommentRepo func(uint64) error
)

func (f *fakeCommentRepo) SaveComment(comment *domain.Comment) (*domain.Comment, map[string]string) {
	return saveCommentRepo(comment)
}
func (f *fakeCommentRepo) GetComment(commentId uint64) (*domain.Comment, error) {
	return getCommentRepo(commentId)
}
func (f *fakeCommentRepo) GetAllComment() ([]domain.Comment, error) {
	return getAllCommentRepo()
}
func (f *fakeCommentRepo) UpdateComment(comment *domain.Comment) (*domain.Comment, map[string]string) {
	return updateCommentRepo(comment)
}
func (f *fakeCommentRepo) DeleteComment(commentId uint64) error {
	return deleteCommentRepo(commentId)
}

//var fakeComment repository.CommentRepository = &fakeCommentRepo{} //this is where the real implementation is swap with our fake implementation
var commentRepoFake application.CommentAppInterface = &fakeCommentRepo{} //this is where the real implementation is swap with our fake implementation

func TestSaveComment_Success(t *testing.T) {
	//Mock the response coming from the infrastructure
	saveCommentRepo = func(user *domain.Comment) (*domain.Comment, map[string]string) {
		return &domain.Comment{
			ID:      1,
			PostID:  1,
			Content: "comment content",
		}, nil
	}
	comment := &domain.Comment{
		ID:      1,
		PostID:  1,
		Content: "comment content",
	}
	c, err := commentRepoFake.SaveComment(comment)
	assert.Nil(t, err)
	assert.EqualValues(t, c.PostID, 1)
	assert.EqualValues(t, c.Content, "comment content")
}

func TestGetComment_Success(t *testing.T) {
	//Mock the response coming from the infrastructure
	getCommentRepo = func(commentId uint64) (*domain.Comment, error) {
		return &domain.Comment{
			ID:      1,
			PostID:  1,
			Content: "comment content",
		}, nil
	}
	commentId := uint64(1)
	c, err := commentRepoFake.GetComment(commentId)
	assert.Nil(t, err)
	assert.EqualValues(t, c.PostID, 1)
	assert.EqualValues(t, c.Content, "comment content")
}

func TestAllComment_Success(t *testing.T) {
	//Mock the response coming from the infrastructure
	getAllCommentRepo = func() ([]domain.Comment, error) {
		return []domain.Comment{
			{
				ID:      1,
				PostID:  1,
				Content: "comment content first",
			},
			{
				ID:      2,
				PostID:  1,
				Content: "comment content second",
			},
		}, nil
	}
	c, err := commentRepoFake.GetAllComment()
	assert.Nil(t, err)
	assert.EqualValues(t, len(c), 2)
}

func TestUpdateComment_Success(t *testing.T) {
	//Mock the response coming from the infrastructure
	updateCommentRepo = func(user *domain.Comment) (*domain.Comment, map[string]string) {
		return &domain.Comment{
			ID:      1,
			PostID:  1,
			Content: "comment content update",
		}, nil
	}
	comment := &domain.Comment{
		ID:      1,
		PostID:  1,
		Content: "comment content update",
	}
	c, err := commentRepoFake.UpdateComment(comment)
	assert.Nil(t, err)
	assert.EqualValues(t, c.PostID, 1)
	assert.EqualValues(t, c.Content, "comment content update")
}

func TestDeleteComment_Success(t *testing.T) {
	//Mock the response coming from the infrastructure
	deleteCommentRepo = func(commentId uint64) error {
		return nil
	}
	commentId := uint64(1)
	err := commentRepoFake.DeleteComment(commentId)
	assert.Nil(t, err)
}
