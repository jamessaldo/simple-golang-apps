package handlers

import (
	"nc-two/application"
	"nc-two/infrastructure/auth"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	PostApp    application.PostAppInterface
	userApp    application.UserAppInterface
	CommentApp application.CommentAppInterface
	tk         auth.TokenInterface
	rd         auth.AuthInterface
}

//Handler constructor
func NewHandler(pApp application.PostAppInterface, cApp application.CommentAppInterface, uApp application.UserAppInterface, rd auth.AuthInterface, tk auth.TokenInterface) *Handler {
	return &Handler{
		PostApp:    pApp,
		userApp:    uApp,
		CommentApp: cApp,
		rd:         rd,
		tk:         tk,
	}
}

type Server struct {
	Handler *Handler
	Router  *gin.Engine
}
