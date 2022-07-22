package handlers

import (
	"nctwo/backend/application"
	"nctwo/backend/infrastructure/auth"
	"nctwo/backend/infrastructure/worker"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	PostApp    application.PostAppInterface
	userApp    application.UserAppInterface
	CommentApp application.CommentAppInterface
	tk         auth.TokenInterface
	rd         auth.AuthInterface
	wk         worker.WorkerInterface
}

//Handler constructor
func NewHandler(pApp application.PostAppInterface, cApp application.CommentAppInterface, uApp application.UserAppInterface, rd auth.AuthInterface, tk auth.TokenInterface, wk worker.WorkerInterface) *Handler {
	return &Handler{
		PostApp:    pApp,
		userApp:    uApp,
		CommentApp: cApp,
		rd:         rd,
		tk:         tk,
		wk:         wk,
	}
}

type Server struct {
	Handler *Handler
	Router  *gin.Engine
}
