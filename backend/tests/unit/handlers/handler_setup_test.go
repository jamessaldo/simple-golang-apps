package handlers

import (
	"nc-two/controllers"
	"nc-two/service_layer"
	"nc-two/service_layer/handlers"
	"nc-two/utils/mock"
)

var (
	userApp   mock.UserAppInterface
	postApp   mock.PostAppInterface
	fakeAuth  mock.AuthInterface
	fakeToken mock.TokenInterface

	s  = handlers.NewUsers(&userApp, &fakeAuth, &fakeToken)          //We use all mocked data here
	f  = handlers.NewPost(&postApp, &userApp, &fakeAuth, &fakeToken) //We use all mocked data here
	au = handlers.NewAuthenticate(&userApp, &fakeAuth, &fakeToken)   //We use all mocked data here

	uow = &service_layer.UnitOfWork{
		Users: &userApp,
		Posts: &postApp,
	}
	bus = controllers.Bootsrap(*uow, &fakeToken, &fakeAuth, controllers.Handler{Users: *s, Posts: *f, Auth: *au})
)
