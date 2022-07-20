package bootstrap

import (
	"nc-two/infrastructure/auth"
	"nc-two/service_layer"
	"nc-two/service_layer/handlers"
)

type Handler struct {
	Users handlers.Users
	Posts handlers.Post
	Auth  handlers.Authenticate
}

type Bootstrap struct {
	UOW     service_layer.UnitOfWork
	TK      auth.TokenInterface
	RD      auth.AuthInterface
	Handler Handler
}

func Bootsrap(uow service_layer.UnitOfWork, tk auth.TokenInterface, rd auth.AuthInterface, handler Handler) *Bootstrap {
	return &Bootstrap{UOW: uow, TK: tk, RD: rd, Handler: handler}
}
