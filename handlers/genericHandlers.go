package handlers

import (

	"fmt"
	
	"net/http"
	
	"strconv"
	
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type GenericHandler[T any] struct {
	db *gorm.DB
}

func NewGenericHandler[T any](db *gorm.DB) *GenericHandler[T] {
	return &GenericHandler[T]{db: db}
}

func (h *GenericHandler[T]) Create(c *gin.Context) {
	var entity T
	if err := c.ShouldBindJSON(&entity); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := h.db.Create(&entity).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create entity"})
		return
	}
	c.JSON(http.StatusOK, entity)
}

func (h *GenericHandler[T]) BulkCreate(c *gin.Context) {
	var entities []T
	if err := c.ShouldBindJSON(&entities); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := h.db.Create(&entities).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create entities"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Bulk data inserted successfully"})
}

func (h *GenericHandler[T]) UploadCSV(c *gin.Context) {
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file upload"})
		return
	}
	defer file.Close()

	records, err := ParseCSV[T](file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse CSV", "details": err.Error()})
		return
	}

	if err := h.db.Create(&records).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert records"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "CSV data uploaded successfully"})
}

func (h *GenericHandler[T]) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var entity T
	if err := h.db.First(&entity, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
		return
	}
	c.JSON(http.StatusOK, entity)
}

func (h *GenericHandler[T]) GetAll(c *gin.Context) {
	var entities []T
	params := c.Request.URL.Query()
	query := h.db.Model(&entities)

	for key, values := range params {
		if len(values) > 0 {
			query = query.Where(fmt.Sprintf("%s = ?", key), values[0])
		}
	}

	if err := query.Find(&entities).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch records"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": entities})
}

func (h *GenericHandler[T]) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	var entity T
	if err := c.ShouldBindJSON(&entity); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := h.db.Model(&entity).Where("id = ?", id).Updates(&entity).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update record"})
		return
	}
	c.JSON(http.StatusOK, entity)
}

func (h *GenericHandler[T]) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var entity T
	if err := h.db.Delete(&entity, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete record"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Record deleted"})
}
