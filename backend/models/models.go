package models

import (
	"time"
	"gorm.io/gorm"
)

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"unique;not null" json:"username"`
	Password  string    `gorm:"not null" json:"-"`
	Token     string    `gorm:"index" json:"token,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Item struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"not null" json:"name"`
	Description string    `json:"description"`
	Price       float64   `gorm:"not null" json:"price"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Cart struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	UserID    uint       `gorm:"not null;index" json:"user_id"`
	User      User       `gorm:"foreignKey:UserID" json:"-"`
	CartItems []CartItem `gorm:"foreignKey:CartID" json:"cart_items,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type CartItem struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CartID    uint      `gorm:"not null;index" json:"cart_id"`
	Cart      Cart      `gorm:"foreignKey:CartID" json:"-"`
	ItemID    uint      `gorm:"not null" json:"item_id"`
	Item      Item      `gorm:"foreignKey:ItemID" json:"item,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Order struct {
	ID         uint        `gorm:"primaryKey" json:"id"`
	UserID     uint        `gorm:"not null;index" json:"user_id"`
	User       User        `gorm:"foreignKey:UserID" json:"-"`
	OrderItems []OrderItem `gorm:"foreignKey:OrderID" json:"order_items,omitempty"`
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at"`
}

type OrderItem struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	OrderID   uint      `gorm:"not null;index" json:"order_id"`
	Order     Order     `gorm:"foreignKey:OrderID" json:"-"`
	ItemID    uint      `gorm:"not null" json:"item_id"`
	Item      Item      `gorm:"foreignKey:ItemID" json:"item,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func MigrateDB(db *gorm.DB) error {
	return db.AutoMigrate(&User{}, &Item{}, &Cart{}, &CartItem{}, &Order{}, &OrderItem{})
}
