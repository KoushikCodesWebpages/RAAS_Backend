// repositories/protected_repo.go
package handlers

import (
	"gorm.io/gorm"
)

type ProtectedRepository[T any] struct {
	db *gorm.DB
}

func NewProtectedRepository[T any](db *gorm.DB) *ProtectedRepository[T] {
	return &ProtectedRepository[T]{db: db}
}

func (r *ProtectedRepository[T]) Create(entity *T) error {
	return r.db.Create(entity).Error
}

func (r *ProtectedRepository[T]) BulkCreate(entities []T) error {
	return r.db.Create(&entities).Error
}

func (r *ProtectedRepository[T]) GetByID(id uint, entity *T) error {
	return r.db.First(entity, id).Error
}

func (r *ProtectedRepository[T]) GetAll(filters map[string]interface{}) ([]T, error) {
	var results []T
	query := r.db.Model(&results)
	for field, value := range filters {
		query = query.Where(field+" = ?", value)
	}
	err := query.Find(&results).Error
	return results, err
}

func (r *ProtectedRepository[T]) Update(id uint, entity *T) error {
	return r.db.Model(entity).Where("id = ?", id).Updates(entity).Error
}

func (r *ProtectedRepository[T]) Delete(id uint, entity *T) error {
	return r.db.Delete(entity, id).Error
}
