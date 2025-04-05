package controllers

import (
	"github.com/gin-gonic/gin"
	"RAAS/repositories"
	"net/http"
	"strconv"
)

// GeneralController - Generic Controller for CRUD Operations
type GeneralController[T any] struct {
	repo *repositories.GeneralRepository[T]
}

// NewGeneralController - Returns a new instance of GeneralController
func NewGeneralController[T any](repo *repositories.GeneralRepository[T]) *GeneralController[T] {
	return &GeneralController[T]{repo: repo}
}

// Create - Handles the creation of a new entity
func (gc *GeneralController[T]) Create(c *gin.Context) {
	var entity T
	if err := c.ShouldBindJSON(&entity); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	createdEntity, err := gc.repo.Create(&entity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create entity"})
		return
	}
	c.JSON(http.StatusOK, createdEntity)
}

// BulkCreate - Handles the creation of multiple entities
func (gc *GeneralController[T]) BulkCreate(c *gin.Context) {
	var entities []T
	if err := c.ShouldBindJSON(&entities); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	err := gc.repo.BulkCreate(&entities)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create entities"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Bulk data inserted successfully"})
}

func (gc *GeneralController[T]) UploadCSV(c *gin.Context) {
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file upload"})
		return
	}
	defer file.Close()

	// Parse CSV data
	records, err := parseCSV[T](file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse CSV"})
		return
	}

	// Perform bulk insert
	err = gc.repo.BulkCreate(&records)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert records"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "CSV data uploaded successfully"})
}

// GetByID - Fetch a record by ID
func (gc *GeneralController[T]) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	entity, err := gc.repo.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
		return
	}

	c.JSON(http.StatusOK, entity)
}

// GetAll - Fetch all records with optional filtering
func (gc *GeneralController[T]) GetAll(c *gin.Context) {
	entities, err := gc.repo.GetAll(c.Request.URL.Query())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch records"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		// "total": total, // Commented out pagination for now
		"data": entities,
	})
}

// Update - Updates a record
func (gc *GeneralController[T]) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var entity T
	if err := c.ShouldBindJSON(&entity); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	updatedEntity, err := gc.repo.Update(uint(id), &entity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update record"})
		return
	}

	c.JSON(http.StatusOK, updatedEntity)
}

// Delete - Removes a record by ID
func (gc *GeneralController[T]) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	err = gc.repo.Delete(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete record"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Record deleted"})
}
