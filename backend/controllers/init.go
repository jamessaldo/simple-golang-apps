package controllers

import (
	"nc-two/bootstrap"
	"nc-two/infrastructure/auth"
	"nc-two/service_layer"
)

type Bootstrap struct {
	UOW     service_layer.UnitOfWork
	TK      auth.TokenInterface
	RD      auth.AuthInterface
	Handler bootstrap.Handler
}
