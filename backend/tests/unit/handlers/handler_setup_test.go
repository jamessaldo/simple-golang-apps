package handlers

import (
	"nc-two/handlers"
	"nc-two/utils/mock"
)

var (
	userApp    mock.UserAppInterface
	PostApp    mock.PostAppInterface
	CommentApp mock.CommentAppInterface
	fakeAuth   mock.AuthInterface
	fakeToken  mock.TokenInterface

	handler = handlers.NewHandler(&PostApp, &CommentApp, &userApp, &fakeAuth, &fakeToken) //We use all mocked data here
)
