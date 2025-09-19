package services

import (
	"errors"

	"github.com/RaihanurRahman2022/PersonalVault/internal/app/entities"
	"github.com/RaihanurRahman2022/PersonalVault/internal/app/repositories"
	"github.com/RaihanurRahman2022/PersonalVault/internal/helper"
)

var (
	ErrInvalidCredentials = errors.New("invalid username or password")
	ErrUserInactive       = errors.New("user is inactive")
	ErrInvalidToken       = errors.New("invalid token")
)

type AuthService interface {
	Login(username, password string) (string, string, error)
	Register(req *entities.RegisterRequest) error
	RefreshToken(refreshToken string) (string, string, error)
}

type AuthServiceImpl struct {
	authRepo repositories.AuthRepository
}

func NewAuthService(authRepo repositories.AuthRepository) AuthService {
	return &AuthServiceImpl{
		authRepo: authRepo,
	}
}

func (r *AuthServiceImpl) Login(username, passwerod string) (string, string, error) {
	user, err := r.authRepo.GetUserByUsername(username)

	if err != nil {
		return "", "", ErrInvalidCredentials
	}

	if !helper.CheckPassword(passwerod, user.Password) {
		return "", "", ErrInvalidCredentials
	}

	accessToken, err := helper.GenerateJWT(user.ID, username)

	if err != nil {
		return "", "", err
	}

	// Generate refresh token (7 days)
	refreshToken, err := helper.GenerateRefreshToken(user.ID)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (r *AuthServiceImpl) RefreshToken(refreshToken string) (string, string, error) {
	claims, err := helper.ValidateRefreshToken(refreshToken)
	if err != nil {
		return "", "", ErrInvalidToken
	}

	user, err := r.authRepo.GetUserById(claims.UserID)
	if err != nil {
		return "", "", ErrInvalidToken
	}

	accessToken, err := helper.GenerateJWT(user.ID, user.UserName)

	if err != nil {
		return "", "", err
	}

	// Generate refresh token (7 days)
	newRefreshToken, err := helper.GenerateRefreshToken(user.ID)
	if err != nil {
		return "", "", err
	}

	return accessToken, newRefreshToken, nil
}

func (r *AuthServiceImpl) Register(req *entities.RegisterRequest) error {
	hashedPassword, err := helper.HashPassword(req.Password)
	if err != nil {
		return err
	}

	user := &entities.User{
		UserName:  req.Username,
		Password:  hashedPassword,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}

	return r.authRepo.Create(user)
}
