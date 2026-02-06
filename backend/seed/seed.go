package main

import (
	"log"
	"shopping-cart-backend/models"
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

	items := []models.Item{
		{Name: "Laptop", Description: "High-performance laptop with 16GB RAM", Price: 999.99},
		{Name: "Smartphone", Description: "Latest smartphone with 5G support", Price: 699.99},
		{Name: "Wireless Headphones", Description: "Noise-cancelling over-ear headphones", Price: 199.99},
		{Name: "Smart Watch", Description: "Fitness tracker with heart rate monitor", Price: 299.99},
		{Name: "Tablet", Description: "10-inch tablet with stylus support", Price: 449.99},
		{Name: "Keyboard", Description: "Mechanical keyboard with RGB lighting", Price: 129.99},
		{Name: "Mouse", Description: "Wireless gaming mouse", Price: 79.99},
		{Name: "Monitor", Description: "27-inch 4K display", Price: 399.99},
	}

	for _, item := range items {
		var existingItem models.Item
		result := db.Where("name = ?", item.Name).First(&existingItem)
		if result.Error == gorm.ErrRecordNotFound {
			db.Create(&item)
			log.Printf("Created item: %s\n", item.Name)
		} else {
			log.Printf("Item already exists: %s\n", item.Name)
		}
	}

	users := []models.User{
		{Username: "admin", Password: "admin"},
		{Username: "user1", Password: "password123"},
	}

	for _, user := range users {
		var existingUser models.User
		result := db.Where("username = ?", user.Username).First(&existingUser)
		if result.Error == gorm.ErrRecordNotFound {
			db.Create(&user)
			log.Printf("Created user: %s\n", user.Username)
		} else {
			log.Printf("User already exists: %s\n", user.Username)
		}
	}

	log.Println("Database seeding completed!")
}
