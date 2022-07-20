package handlers

import (
	"nc-two/service_layer/handlers"
	"nc-two/utils/mock"
)

var (
	userApp   mock.UserAppInterface
	postApp   mock.PostAppInterface
	fakeAuth  mock.AuthInterface
	fakeToken mock.TokenInterface

	handler = handlers.NewHandler(&postApp, &userApp, &fakeAuth, &fakeToken) //We use all mocked data here
)
