package mock

import (
	"mime/multipart"
	"nc-two/domain"
	"nc-two/infrastructure/auth"
	"net/http"
)

//UserAppInterface is a mock user app interface
type UserAppInterface struct {
	SaveUserFn                  func(*domain.User) (*domain.User, map[string]string)
	GetUsersFn                  func() ([]domain.User, error)
	GetUserFn                   func(uint64) (*domain.User, error)
	GetUserByEmailAndPasswordFn func(*domain.User) (*domain.User, map[string]string)
}

//SaveUser calls the SaveUserFn
func (u *UserAppInterface) SaveUser(user *domain.User) (*domain.User, map[string]string) {
	return u.SaveUserFn(user)
}

//GetUsersFn calls the GetUsers
func (u *UserAppInterface) GetUsers() ([]domain.User, error) {
	return u.GetUsersFn()
}

//GetUserFn calls the GetUser
func (u *UserAppInterface) GetUser(userId uint64) (*domain.User, error) {
	return u.GetUserFn(userId)
}

//GetUserByEmailAndPasswordFn calls the GetUserByEmailAndPassword
func (u *UserAppInterface) GetUserByEmailAndPassword(user *domain.User) (*domain.User, map[string]string) {
	return u.GetUserByEmailAndPasswordFn(user)
}

//PostAppInterface is a mock post app interface
type PostAppInterface struct {
	SavePostFn   func(*domain.Post) (*domain.Post, map[string]string)
	GetAllPostFn func() ([]domain.Post, error)
	GetPostFn    func(uint64) (*domain.Post, error)
	UpdatePostFn func(*domain.Post) (*domain.Post, map[string]string)
	DeletePostFn func(uint64) error
}

func (p *PostAppInterface) SavePost(post *domain.Post) (*domain.Post, map[string]string) {
	return p.SavePostFn(post)
}
func (p *PostAppInterface) GetAllPost() ([]domain.Post, error) {
	return p.GetAllPostFn()
}
func (p *PostAppInterface) GetPost(postId uint64) (*domain.Post, error) {
	return p.GetPostFn(postId)
}
func (p *PostAppInterface) UpdatePost(post *domain.Post) (*domain.Post, map[string]string) {
	return p.UpdatePostFn(post)
}
func (p *PostAppInterface) DeletePost(postId uint64) error {
	return p.DeletePostFn(postId)
}

//CommentAppInterface is a mock post app interface
type CommentAppInterface struct {
	SaveCommentFn   func(*domain.Comment) (*domain.Comment, map[string]string)
	GetAllCommentFn func() ([]domain.Comment, error)
	GetCommentFn    func(uint64) (*domain.Comment, error)
	UpdateCommentFn func(*domain.Comment) (*domain.Comment, map[string]string)
	DeleteCommentFn func(uint64) error
}

func (c *CommentAppInterface) SaveComment(post *domain.Comment) (*domain.Comment, map[string]string) {
	return c.SaveCommentFn(post)
}
func (c *CommentAppInterface) GetAllComment() ([]domain.Comment, error) {
	return c.GetAllCommentFn()
}
func (c *CommentAppInterface) GetComment(postId uint64) (*domain.Comment, error) {
	return c.GetCommentFn(postId)
}
func (c *CommentAppInterface) UpdateComment(post *domain.Comment) (*domain.Comment, map[string]string) {
	return c.UpdateCommentFn(post)
}
func (c *CommentAppInterface) DeleteComment(postId uint64) error {
	return c.DeleteCommentFn(postId)
}

//AuthInterface is a mock auth interface
type AuthInterface struct {
	CreateAuthFn    func(uint64, *auth.TokenDetails) error
	FetchAuthFn     func(string) (uint64, error)
	DeleteRefreshFn func(string) error
	DeleteTokensFn  func(*auth.AccessDetails) error
}

func (a *AuthInterface) DeleteRefresh(refreshUuid string) error {
	return a.DeleteRefreshFn(refreshUuid)
}
func (a *AuthInterface) DeleteTokens(authD *auth.AccessDetails) error {
	return a.DeleteTokensFn(authD)
}
func (a *AuthInterface) FetchAuth(uuid string) (uint64, error) {
	return a.FetchAuthFn(uuid)
}
func (a *AuthInterface) CreateAuth(userId uint64, authD *auth.TokenDetails) error {
	return a.CreateAuthFn(userId, authD)
}

//TokenInterface is a mock token interface
type TokenInterface struct {
	CreateTokenFn          func(userId uint64) (*auth.TokenDetails, error)
	ExtractTokenMetadataFn func(*http.Request) (*auth.AccessDetails, error)
}

func (t *TokenInterface) CreateToken(userid uint64) (*auth.TokenDetails, error) {
	return t.CreateTokenFn(userid)
}
func (t *TokenInterface) ExtractTokenMetadata(r *http.Request) (*auth.AccessDetails, error) {
	return t.ExtractTokenMetadataFn(r)
}

type UploadFileInterface struct {
	UploadFileFn func(file *multipart.FileHeader) (string, error)
}

func (up *UploadFileInterface) UploadFile(file *multipart.FileHeader) (string, error) {
	return up.UploadFileFn(file)
}
