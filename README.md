# Shopping Cart Application

A full-stack shopping cart web application with Go backend and React frontend.

## Backend (Go)

### Technology Stack
- **Gin** - Web framework for REST APIs
- **GORM** - ORM for database operations
- **Ginkgo** - BDD testing framework
- **SQLite** - Database

### Models
- User
- Item
- Cart
- CartItem
- Order
- OrderItem

### API Endpoints

#### User Management
- `POST /users` - Create a new user (signup)
- `GET /users` - Get all users
- `POST /users/login` - User login (returns token)

#### Items
- `POST /items` - Create a new item
- `GET /items` - Get all items

#### Cart (Requires Authentication)
- `POST /carts` - Add item to cart
- `GET /carts` - Get current user's cart

#### Orders (Requires Authentication)
- `POST /orders` - Checkout (convert cart to order)
- `GET /orders` - Get user's order history

### Running the Backend

```bash
cd backend
go mod download
go run main.go
```

The backend will start on `http://localhost:8080`

### Running Tests

```bash
cd backend
go test ./tests -v
```

## Frontend (React)

### Technology Stack
- **React** - UI library
- **Vite** - Build tool and dev server

### Features
- Login screen with username/password authentication
- Items list with click-to-add-to-cart functionality
- Checkout button to convert cart to order
- Cart button to view cart contents
- Order History button to view past orders

### Running the Frontend

```bash
npm install
npm run dev
```

The frontend will start on `http://localhost:3000`

### Building for Production

```bash
npm run build
```

## Application Flow

1. User logs in with username and password
2. On successful login, receives a token (stored in localStorage)
3. Views list of available items
4. Clicks on items to add them to cart
5. Can view cart contents using the Cart button
6. Uses Checkout button to convert cart to an order
7. Can view order history using Order History button

## Business Rules

- Each user can have only one active token at a time
- Each user can have only one cart
- Cart is converted to an order on checkout
- No inventory or quantity management
- Authentication required for cart and order operations

## Sample Usage

### Create a User
```bash
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"testpass"}'
```

### Login
```bash
curl -X POST http://localhost:8080/users/login \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"testpass"}'
```

### Create Items
```bash
curl -X POST http://localhost:8080/items \
  -H "Content-Type: application/json" \
  -d '{"name":"Laptop","description":"High-performance laptop","price":999.99}'
```

### Add to Cart
```bash
curl -X POST http://localhost:8080/carts \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{"item_id":1}'
```

### Checkout
```bash
curl -X POST http://localhost:8080/orders \
  -H "Authorization: Bearer YOUR_TOKEN"
```