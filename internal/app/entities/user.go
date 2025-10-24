package entities

import "github.com/google/uuid"

type User struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()" example:"550e8400-e29b-41d4-a716-446655440000"`
	UserName  string    `json:"username" gorm:"uniqueIndex;not null" example:"john_doe"`
	Password  string    `json:"-" gorm:"not null"`
	FirstName string    `json:"first_name" gorm:"not null" example:"John"`
	LastName  string    `json:"last_name" gorm:"not null" example:"Doe"`
}
