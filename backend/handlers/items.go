package handlers

import (
	"net/http"
	"shopping-cart-backend/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ItemHandler struct {
	DB *gorm.DB
}

type CreateItemRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required"`
}

func (h *ItemHandler) CreateItem(c *gin.Context) {
	var req CreateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item := models.Item{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
	}

	if err := h.DB.Create(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create item"})
		return
	}

	c.JSON(http.StatusCreated, item)
}

func (h *ItemHandler) GetItems(c *gin.Context) {
	var items []models.Item
	h.DB.Find(&items)
	c.JSON(http.StatusOK, items)
}
