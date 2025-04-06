// handlers/protected_handler.go
package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProtectedHandler[T any] struct {
	repo *ProtectedRepository[T]
}

func NewProtectedHandler[T any](db *gorm.DB) *ProtectedHandler[T] {
	return &ProtectedHandler[T]{
		repo: NewProtectedRepository[T](db),
	}
}

func (h *ProtectedHandler[T]) CreateProtected(c *gin.Context) {
	var entity T
	if err := c.ShouldBindJSON(&entity); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}
	if err := h.repo.Create(&entity); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Create failed"})
		return
	}
	c.JSON(http.StatusOK, entity)
}

func (h *ProtectedHandler[T]) BulkCreateProtected(c *gin.Context) {
	var entities []T
	if err := c.ShouldBindJSON(&entities); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}
	if err := h.repo.BulkCreate(entities); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Bulk create failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Bulk create successful"})
}

func (h *ProtectedHandler[T]) GetByIDProtected(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var entity T
	if err := h.repo.GetByID(uint(id), &entity); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
		return
	}
	c.JSON(http.StatusOK, entity)
}

func (h *ProtectedHandler[T]) GetAllProtected(c *gin.Context) {
	allowedFilters := []string{"status", "type", "category"} // customize this
	filters := map[string]interface{}{}
	for _, f := range allowedFilters {
		if val := c.Query(f); val != "" {
			filters[f] = val
		}
	}
	entities, err := h.repo.GetAll(filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fetch failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": entities})
}

func (h *ProtectedHandler[T]) UpdateProtected(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var entity T
	if err := c.ShouldBindJSON(&entity); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid body"})
		return
	}
	if err := h.repo.Update(uint(id), &entity); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Update failed"})
		return
	}
	c.JSON(http.StatusOK, entity)
}

func (h *ProtectedHandler[T]) DeleteProtected(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var entity T
	if err := h.repo.Delete(uint(id), &entity); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Delete failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Deleted"})
}
