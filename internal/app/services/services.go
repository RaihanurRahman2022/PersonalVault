package services

import "github.com/RaihanurRahman2022/PersonalVault/internal/app/repositories"

type Services struct {
	User UserService
	Auth AuthService
}

func NewServices(repo *repositories.Repositories) *Services {
	return &Services{
		User: NewUserService(repo.User),
		Auth: NewAuthService(repo.Auth),
	}
}
