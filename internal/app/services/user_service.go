package services

import (
	"github.com/RaihanurRahman2022/PersonalVault/internal/app/entities"
	"github.com/RaihanurRahman2022/PersonalVault/internal/app/repositories"
	"github.com/google/uuid"
)

type UserService interface {
	GetUserByID(id uuid.UUID) (*entities.User, error)
}

type UserServiceImpl struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) UserService {
	return &UserServiceImpl{
		userRepo: userRepo,
	}
}

func (s *UserServiceImpl) GetUserByID(id uuid.UUID) (*entities.User, error) {
	return s.userRepo.GetByID(id)
}
