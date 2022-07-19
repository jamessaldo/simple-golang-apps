package application

import (
	"nc-two/domain/entity"
	"testing"

	"github.com/stretchr/testify/assert"
)

//IF YOU HAVE TIME, YOU CAN TEST ALL THE METHODS FAILURES

type fakePostRepo struct{}

var (
	savePostRepo   func(*entity.Post) (*entity.Post, map[string]string)
	getPostRepo    func(uint64) (*entity.Post, error)
	getAllPostRepo func() ([]entity.Post, error)
	updatePostRepo func(*entity.Post) (*entity.Post, map[string]string)
	deletePostRepo func(uint64) error
)

func (f *fakePostRepo) SavePost(post *entity.Post) (*entity.Post, map[string]string) {
	return savePostRepo(post)
}
func (f *fakePostRepo) GetPost(postId uint64) (*entity.Post, error) {
	return getPostRepo(postId)
}
func (f *fakePostRepo) GetAllPost() ([]entity.Post, error) {
	return getAllPostRepo()
}
func (f *fakePostRepo) UpdatePost(post *entity.Post) (*entity.Post, map[string]string) {
	return updatePostRepo(post)
}
func (f *fakePostRepo) DeletePost(postId uint64) error {
	return deletePostRepo(postId)
}

//var fakePost repository.PostRepository = &fakePostRepo{} //this is where the real implementation is swap with our fake implementation
var postAppFake PostAppInterface = &fakePostRepo{} //this is where the real implementation is swap with our fake implementation

func TestSavePost_Success(t *testing.T) {
	//Mock the response coming from the infrastructure
	savePostRepo = func(user *entity.Post) (*entity.Post, map[string]string) {
		return &entity.Post{
			ID:          1,
			Title:       "post title",
			Description: "post description",
			UserID:      1,
		}, nil
	}
	post := &entity.Post{
		ID:          1,
		Title:       "post title",
		Description: "post description",
		UserID:      1,
	}
	f, err := postAppFake.SavePost(post)
	assert.Nil(t, err)
	assert.EqualValues(t, f.Title, "post title")
	assert.EqualValues(t, f.Description, "post description")
	assert.EqualValues(t, f.UserID, 1)
}

func TestGetPost_Success(t *testing.T) {
	//Mock the response coming from the infrastructure
	getPostRepo = func(postId uint64) (*entity.Post, error) {
		return &entity.Post{
			ID:          1,
			Title:       "post title",
			Description: "post description",
			UserID:      1,
		}, nil
	}
	postId := uint64(1)
	f, err := postAppFake.GetPost(postId)
	assert.Nil(t, err)
	assert.EqualValues(t, f.Title, "post title")
	assert.EqualValues(t, f.Description, "post description")
	assert.EqualValues(t, f.UserID, 1)
}

func TestAllPost_Success(t *testing.T) {
	//Mock the response coming from the infrastructure
	getAllPostRepo = func() ([]entity.Post, error) {
		return []entity.Post{
			{
				ID:          1,
				Title:       "post title first",
				Description: "post description first",
				UserID:      1,
			},
			{
				ID:          2,
				Title:       "post title second",
				Description: "post description second",
				UserID:      1,
			},
		}, nil
	}
	f, err := postAppFake.GetAllPost()
	assert.Nil(t, err)
	assert.EqualValues(t, len(f), 2)
}

func TestUpdatePost_Success(t *testing.T) {
	//Mock the response coming from the infrastructure
	updatePostRepo = func(user *entity.Post) (*entity.Post, map[string]string) {
		return &entity.Post{
			ID:          1,
			Title:       "post title update",
			Description: "post description update",
			UserID:      1,
		}, nil
	}
	post := &entity.Post{
		ID:          1,
		Title:       "post title update",
		Description: "post description update",
		UserID:      1,
	}
	f, err := postAppFake.UpdatePost(post)
	assert.Nil(t, err)
	assert.EqualValues(t, f.Title, "post title update")
	assert.EqualValues(t, f.Description, "post description update")
	assert.EqualValues(t, f.UserID, 1)
}

func TestDeletePost_Success(t *testing.T) {
	//Mock the response coming from the infrastructure
	deletePostRepo = func(postId uint64) error {
		return nil
	}
	postId := uint64(1)
	err := postAppFake.DeletePost(postId)
	assert.Nil(t, err)
}
