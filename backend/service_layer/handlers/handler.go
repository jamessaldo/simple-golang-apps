package handlers

import (
	"nc-two/application"
	"nc-two/infrastructure/auth"
)

type Handler struct {
	postApp application.PostAppInterface
	userApp application.UserAppInterface
	tk      auth.TokenInterface
	rd      auth.AuthInterface
}

//Handler constructor
func NewHandler(fApp application.PostAppInterface, uApp application.UserAppInterface, rd auth.AuthInterface, tk auth.TokenInterface) *Handler {
	return &Handler{
		postApp: fApp,
		userApp: uApp,
		rd:      rd,
		tk:      tk,
	}
}
