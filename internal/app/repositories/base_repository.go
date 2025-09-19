package repositories

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BaseRepository[T any] interface {
	GetByID(id uuid.UUID) (*T, error)
	Create(entity *T) error
	Update(entity *T) error
	Delete(id uuid.UUID) error
	List(page, pageSize int) ([]T, int64, error)
}

type BaseRepositoryImpl[T any] struct {
	db *gorm.DB
}

func (r *BaseRepositoryImpl[T]) GetByID(id uuid.UUID) (*T, error) {
	var entity T
	err := r.db.Where("id = ?", id).First(&entity).Error
	if err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *BaseRepositoryImpl[T]) Create(entity *T) error {
	return r.db.Create(entity).Error
}

func (r *BaseRepositoryImpl[T]) Update(entity *T) error {
	return r.db.Save(entity).Error
}

func (r *BaseRepositoryImpl[T]) Delete(id uuid.UUID) error {
	return r.db.Delete(new(T), id).Error
}

func (r *BaseRepositoryImpl[T]) List(page, pageSize int) ([]T, int64, error) {
	var entities []T
	var total int64

	query := r.db.Model(new(T))

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize

	if err := query.Offset(offset).Limit(pageSize).Find(&entities).Error; err != nil {
		return nil, 0, err
	}

	return entities, total, nil
}
