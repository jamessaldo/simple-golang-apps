package mock

import (
	"mime/multipart"
	"nc-two/domain/entity"
	"nc-two/infrastructure/auth"
	"net/http"
)

//UserAppInterface is a mock user app interface
type UserAppInterface struct {
	SaveUserFn                  func(*entity.User) (*entity.User, map[string]string)
	GetUsersFn                  func() ([]entity.User, error)
	GetUserFn                   func(uint64) (*entity.User, error)
	GetUserByEmailAndPasswordFn func(*entity.User) (*entity.User, map[string]string)
}

//SaveUser calls the SaveUserFn
func (u *UserAppInterface) SaveUser(user *entity.User) (*entity.User, map[string]string) {
	return u.SaveUserFn(user)
}

//GetUsersFn calls the GetUsers
func (u *UserAppInterface) GetUsers() ([]entity.User, error) {
	return u.GetUsersFn()
}

//GetUserFn calls the GetUser
func (u *UserAppInterface) GetUser(userId uint64) (*entity.User, error) {
	return u.GetUserFn(userId)
}

//GetUserByEmailAndPasswordFn calls the GetUserByEmailAndPassword
func (u *UserAppInterface) GetUserByEmailAndPassword(user *entity.User) (*entity.User, map[string]string) {
	return u.GetUserByEmailAndPasswordFn(user)
}

//PostAppInterface is a mock post app interface
type PostAppInterface struct {
	SavePostFn   func(*entity.Post) (*entity.Post, map[string]string)
	GetAllPostFn func() ([]entity.Post, error)
	GetPostFn    func(uint64) (*entity.Post, error)
	UpdatePostFn func(*entity.Post) (*entity.Post, map[string]string)
	DeletePostFn func(uint64) error
}

func (f *PostAppInterface) SavePost(post *entity.Post) (*entity.Post, map[string]string) {
	return f.SavePostFn(post)
}
func (f *PostAppInterface) GetAllPost() ([]entity.Post, error) {
	return f.GetAllPostFn()
}
func (f *PostAppInterface) GetPost(postId uint64) (*entity.Post, error) {
	return f.GetPostFn(postId)
}
func (f *PostAppInterface) UpdatePost(post *entity.Post) (*entity.Post, map[string]string) {
	return f.UpdatePostFn(post)
}
func (f *PostAppInterface) DeletePost(postId uint64) error {
	return f.DeletePostFn(postId)
}

//AuthInterface is a mock auth interface
type AuthInterface struct {
	CreateAuthFn    func(uint64, *auth.TokenDetails) error
	FetchAuthFn     func(string) (uint64, error)
	DeleteRefreshFn func(string) error
	DeleteTokensFn  func(*auth.AccessDetails) error
}

func (f *AuthInterface) DeleteRefresh(refreshUuid string) error {
	return f.DeleteRefreshFn(refreshUuid)
}
func (f *AuthInterface) DeleteTokens(authD *auth.AccessDetails) error {
	return f.DeleteTokensFn(authD)
}
func (f *AuthInterface) FetchAuth(uuid string) (uint64, error) {
	return f.FetchAuthFn(uuid)
}
func (f *AuthInterface) CreateAuth(userId uint64, authD *auth.TokenDetails) error {
	return f.CreateAuthFn(userId, authD)
}

//TokenInterface is a mock token interface
type TokenInterface struct {
	CreateTokenFn          func(userId uint64) (*auth.TokenDetails, error)
	ExtractTokenMetadataFn func(*http.Request) (*auth.AccessDetails, error)
}

func (f *TokenInterface) CreateToken(userid uint64) (*auth.TokenDetails, error) {
	return f.CreateTokenFn(userid)
}
func (f *TokenInterface) ExtractTokenMetadata(r *http.Request) (*auth.AccessDetails, error) {
	return f.ExtractTokenMetadataFn(r)
}

type UploadFileInterface struct {
	UploadFileFn func(file *multipart.FileHeader) (string, error)
}

func (up *UploadFileInterface) UploadFile(file *multipart.FileHeader) (string, error) {
	return up.UploadFileFn(file)
}
