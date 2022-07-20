package handlers

import (
	"nc-two/adapters"
	"nc-two/domain/models"
	"nc-two/infrastructure/auth"
)

type Post struct {
	postApp adapters.PostRepository
	userApp adapters.UserRepository
	tk      auth.TokenInterface
	rd      auth.AuthInterface
}

//Post constructor
func NewPost(fRepo adapters.PostRepository, uRepo adapters.UserRepository, rd auth.AuthInterface, tk auth.TokenInterface) *Post {
	return &Post{
		postApp: fRepo,
		userApp: uRepo,
		rd:      rd,
		tk:      tk,
	}
}

func (fo *Post) SavePost(post *models.Post) (*models.Post, map[string]string) {
	savedPost, saveErr := fo.postApp.SavePost(post)
	if saveErr != nil {
		return nil, saveErr
	}
	return savedPost, nil
}

func (fo *Post) UpdatePost(post *models.Post) (*models.Post, map[string]string) {
	updatedPost, dbUpdateErr := fo.postApp.UpdatePost(post)
	if dbUpdateErr != nil {
		return nil, dbUpdateErr
	}

	return updatedPost, nil
}

func (fo *Post) GetAllPost() ([]models.Post, error) {
	allpost, err := fo.postApp.GetAllPost()
	if err != nil {
		return nil, err
	}
	return allpost, nil
}

func (fo *Post) GetPostAndCreator(postId uint64) (map[string]interface{}, error) {
	post, err := fo.postApp.GetPost(postId)
	if err != nil {
		return nil, err
	}
	user, err := fo.userApp.GetUser(post.UserID)
	if err != nil {
		return nil, err
	}
	postAndUser := map[string]interface{}{
		"post":    post,
		"creator": user.PublicUser(),
	}
	return postAndUser, nil
}

func (fo *Post) DeletePost(postId uint64) (string, error) {
	err := fo.postApp.DeletePost(postId)
	if err != nil {
		return "", err
	}
	return "Post deleted successfully", nil
}
