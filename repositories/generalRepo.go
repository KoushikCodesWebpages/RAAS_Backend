package repositories

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"net/url"
)

// GeneralRepository - Generic Repository for CRUD Operations
type GeneralRepository[T any] struct {
	db *gorm.DB
}

// NewGeneralRepository - Returns a new instance of GeneralRepository
func NewGeneralRepository[T any](db *gorm.DB) *GeneralRepository[T] {
	return &GeneralRepository[T]{db: db}
}

// Create - Adds a new record
func (r *GeneralRepository[T]) Create(entity *T) (*T, error) {
	if err := r.db.Create(entity).Error; err != nil {
		return nil, err
	}
	return entity, nil
}

// BulkCreate - Adds multiple records at once
func (r *GeneralRepository[T]) BulkCreate(entities *[]T) error {
	if err := r.db.Create(entities).Error; err != nil {
		return err
	}
	return nil
}

// GetByID - Fetch a record by ID
func (r *GeneralRepository[T]) GetByID(id uint) (*T, error) {
	var entity T
	if err := r.db.First(&entity, id).Error; err != nil {
		return nil, errors.New("record not found")
	}
	return &entity, nil
}

// GetAll - Fetch all records with optional filtering
func (r *GeneralRepository[T]) GetAll(queryParams url.Values) ([]T, error) {
	var entities []T
	query := r.db.Model(&entities)

	// Filtering based on query parameters
	for key, values := range queryParams {
		if len(values) > 0 {
			query = query.Where(fmt.Sprintf("%s = ?", key), values[0])
		}
	}

	// Pagination disabled for now
	// page, _ := strconv.Atoi(queryParams.Get("page"))
	// pageSize, _ := strconv.Atoi(queryParams.Get("page_size"))
	// if page == 0 { page = 1 }
	// if pageSize == 0 { pageSize = 10 }
	// offset := (page - 1) * pageSize
	// query.Offset(offset).Limit(pageSize)

	query.Find(&entities)
	return entities, nil
}

// Update - Updates an existing record
func (r *GeneralRepository[T]) Update(id uint, entity *T) (*T, error) {
	if err := r.db.Model(&entity).Where("id = ?", id).Updates(entity).Error; err != nil {
		return nil, err
	}
	return entity, nil
}

// Delete - Removes a record by ID
func (r *GeneralRepository[T]) Delete(id uint) error {
	var entity T
	if err := r.db.Delete(&entity, id).Error; err != nil {
		return errors.New("failed to delete record")
	}
	return nil
}
