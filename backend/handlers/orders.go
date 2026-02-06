package handlers

import (
	"net/http"
	"shopping-cart-backend/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type OrderHandler struct {
	DB *gorm.DB
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	userID := c.GetUint("user_id")

	var cart models.Cart
	err := h.DB.Where("user_id = ?", userID).Preload("CartItems").First(&cart).Error
	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cart is empty"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get cart"})
		return
	}

	if len(cart.CartItems) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cart is empty"})
		return
	}

	order := models.Order{UserID: userID}
	if err := h.DB.Create(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
		return
	}

	for _, cartItem := range cart.CartItems {
		orderItem := models.OrderItem{
			OrderID: order.ID,
			ItemID:  cartItem.ItemID,
		}
		h.DB.Create(&orderItem)
	}

	h.DB.Where("cart_id = ?", cart.ID).Delete(&models.CartItem{})
	h.DB.Delete(&cart)

	h.DB.Preload("OrderItems.Item").First(&order, order.ID)
	c.JSON(http.StatusCreated, order)
}

func (h *OrderHandler) GetOrders(c *gin.Context) {
	userID := c.GetUint("user_id")

	var orders []models.Order
	h.DB.Where("user_id = ?", userID).Preload("OrderItems.Item").Find(&orders)
	c.JSON(http.StatusOK, orders)
}
