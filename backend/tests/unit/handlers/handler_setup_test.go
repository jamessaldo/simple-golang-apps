package handlers

import (
	"nctwo/backend/handlers"
	"nctwo/backend/utils/mock"
)

var (
	userApp    mock.UserAppInterface
	PostApp    mock.PostAppInterface
	CommentApp mock.CommentAppInterface
	fakeAuth   mock.AuthInterface
	fakeToken  mock.TokenInterface
	fakeWorker mock.WorkerInterface

	handler = handlers.NewHandler(&PostApp, &CommentApp, &userApp, &fakeAuth, &fakeToken, &fakeWorker) //We use all mocked data here
)
