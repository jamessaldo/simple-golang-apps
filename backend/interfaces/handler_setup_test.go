package interfaces

import "nc-two/utils/mock"

var (
	userApp    mock.UserAppInterface
	postApp    mock.PostAppInterface
	fakeUpload mock.UploadFileInterface
	fakeAuth   mock.AuthInterface
	fakeToken  mock.TokenInterface

	s  = NewUsers(&userApp, &fakeAuth, &fakeToken)                       //We use all mocked data here
	f  = NewPost(&postApp, &userApp, &fakeUpload, &fakeAuth, &fakeToken) //We use all mocked data here
	au = NewAuthenticate(&userApp, &fakeAuth, &fakeToken)                //We use all mocked data here

)
