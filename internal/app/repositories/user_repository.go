package repositories

import (
	"github.com/RaihanurRahman2022/PersonalVault/internal/app/entities"
	"gorm.io/gorm"
)

type UserRepository interface {
	BaseRepository[entities.User]
}

type UserRepositoryImpl struct {
	BaseRepositoryImpl[entities.User]
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &UserRepositoryImpl{
		BaseRepositoryImpl: BaseRepositoryImpl[entities.User]{db: db},
	}
}
