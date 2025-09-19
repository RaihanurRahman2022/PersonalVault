package repositories

import (
	"github.com/RaihanurRahman2022/PersonalVault/internal/app/entities"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuthRepository interface {
	GetUserByUsername(username string) (*entities.User, error)
	GetUserById(id uuid.UUID) (*entities.User, error)
	Create(user *entities.User) error
}

type AuthRepositoryImpl struct {
	db *gorm.DB
}

func NewAuthReporsitory(db *gorm.DB) AuthRepository {
	return &AuthRepositoryImpl{
		db: db,
	}
}

func (r *AuthRepositoryImpl) GetUserByUsername(username string) (*entities.User, error) {
	var user entities.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *AuthRepositoryImpl) GetUserById(id uuid.UUID) (*entities.User, error) {
	var user entities.User
	err := r.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *AuthRepositoryImpl) Create(user *entities.User) error {
	return r.db.Create(user).Error
}
