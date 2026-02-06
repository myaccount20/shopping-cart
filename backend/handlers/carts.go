package handlers

import (
	"net/http"
	"shopping-cart-backend/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CartHandler struct {
	DB *gorm.DB
}

type AddToCartRequest struct {
	ItemID uint `json:"item_id" binding:"required"`
}

func (h *CartHandler) AddToCart(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req AddToCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var item models.Item
	if err := h.DB.First(&item, req.ItemID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		return
	}

	var cart models.Cart
	err := h.DB.Where("user_id = ?", userID).First(&cart).Error
	if err == gorm.ErrRecordNotFound {
		cart = models.Cart{UserID: userID}
		if err := h.DB.Create(&cart).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create cart"})
			return
		}
	}

	cartItem := models.CartItem{
		CartID: cart.ID,
		ItemID: req.ItemID,
	}

	if err := h.DB.Create(&cartItem).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add item to cart"})
		return
	}

	c.JSON(http.StatusCreated, cartItem)
}

func (h *CartHandler) GetCart(c *gin.Context) {
	userID := c.GetUint("user_id")

	var cart models.Cart
	err := h.DB.Where("user_id = ?", userID).Preload("CartItems.Item").First(&cart).Error
	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusOK, gin.H{"cart_items": []models.CartItem{}})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get cart"})
		return
	}

	c.JSON(http.StatusOK, cart)
}
