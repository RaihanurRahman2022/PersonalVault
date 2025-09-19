package repositories

import "gorm.io/gorm"

type Repositories struct {
	User UserRepository
	Auth AuthRepository
}

func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		User: NewUserRepository(db),
		Auth: NewAuthReporsitory(db),
	}
}
