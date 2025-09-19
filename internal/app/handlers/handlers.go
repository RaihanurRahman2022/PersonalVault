package handlers

import "github.com/RaihanurRahman2022/PersonalVault/internal/app/services"

type Handlers struct {
	UserHandler *UserHandler
}

func NewHandlers(srvc *services.Services) *Handlers {
	return &Handlers{
		UserHandler: NewUserhandler(srvc.User),
	}
}
