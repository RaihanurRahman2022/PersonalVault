package repositories

import "gorm.io/gorm"

type Repositories struct {
	User   UserRepository
	Auth   AuthRepository
	Driver DriverRepository
}

func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		User:   NewUserRepository(db),
		Auth:   NewAuthReporsitory(db),
		Driver: NewDriverRepository(),
	}
}
