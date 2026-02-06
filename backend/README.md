# Shopping Cart Backend

Go backend for the Shopping Cart application using Gin and GORM.

## Prerequisites

- Go 1.21 or higher

## Setup

1. Install dependencies:
```bash
go mod download
```

2. Seed the database with sample data (optional):
```bash
cd seed
go run seed.go
cd ..
```

This will create:
- Sample items (Laptop, Smartphone, Headphones, etc.)
- Test users (admin/admin, user1/password123)

3. Run the server:
```bash
go run main.go
```

The server will start on `http://localhost:8080`

## Running Tests

```bash
go test ./tests -v
```

## API Endpoints

### Public Endpoints
- `POST /users` - Create a new user
- `GET /users` - Get all users
- `POST /users/login` - User login
- `POST /items` - Create a new item
- `GET /items` - Get all items

### Protected Endpoints (Require Authorization header)
- `POST /carts` - Add item to cart
- `GET /carts` - Get cart
- `POST /orders` - Create order (checkout)
- `GET /orders` - Get orders

## Database

Uses SQLite with the database file `shopping_cart.db` in the backend directory.
