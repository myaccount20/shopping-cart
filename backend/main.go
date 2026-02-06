package main

import (
	"log"
	"shopping-cart-backend/handlers"
	"shopping-cart-backend/middleware"
	"shopping-cart-backend/models"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(sqlite.Open("shopping_cart.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err := models.MigrateDB(db); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	authHandler := &handlers.AuthHandler{DB: db}
	itemHandler := &handlers.ItemHandler{DB: db}
	cartHandler := &handlers.CartHandler{DB: db}
	orderHandler := &handlers.OrderHandler{DB: db}

	r.POST("/users", authHandler.Signup)
	r.GET("/users", authHandler.GetUsers)
	r.POST("/users/login", authHandler.Login)

	r.POST("/items", itemHandler.CreateItem)
	r.GET("/items", itemHandler.GetItems)

	authenticated := r.Group("/")
	authenticated.Use(middleware.AuthMiddleware(db))
	{
		authenticated.POST("/carts", cartHandler.AddToCart)
		authenticated.GET("/carts", cartHandler.GetCart)
		authenticated.POST("/orders", orderHandler.CreateOrder)
		authenticated.GET("/orders", orderHandler.GetOrders)
	}

	log.Println("Server starting on :8080")
	r.Run(":8080")
}
