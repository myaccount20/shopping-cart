package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"shopping-cart-backend/handlers"
	"shopping-cart-backend/middleware"
	"shopping-cart-backend/models"
	"testing"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestAPI(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "API Suite")
}

var _ = Describe("Shopping Cart API", func() {
	var (
		db     *gorm.DB
		router *gin.Engine
	)

	BeforeEach(func() {
		var err error
		db, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		Expect(err).NotTo(HaveOccurred())

		err = models.MigrateDB(db)
		Expect(err).NotTo(HaveOccurred())

		gin.SetMode(gin.TestMode)
		router = gin.Default()

		authHandler := &handlers.AuthHandler{DB: db}
		itemHandler := &handlers.ItemHandler{DB: db}
		cartHandler := &handlers.CartHandler{DB: db}
		orderHandler := &handlers.OrderHandler{DB: db}

		router.POST("/users", authHandler.Signup)
		router.GET("/users", authHandler.GetUsers)
		router.POST("/users/login", authHandler.Login)
		router.POST("/items", itemHandler.CreateItem)
		router.GET("/items", itemHandler.GetItems)

		authenticated := router.Group("/")
		authenticated.Use(middleware.AuthMiddleware(db))
		{
			authenticated.POST("/carts", cartHandler.AddToCart)
			authenticated.GET("/carts", cartHandler.GetCart)
			authenticated.POST("/orders", orderHandler.CreateOrder)
			authenticated.GET("/orders", orderHandler.GetOrders)
		}
	})

	Describe("User Registration and Login", func() {
		It("should register a new user", func() {
			payload := map[string]string{"username": "testuser", "password": "testpass"}
			body, _ := json.Marshal(payload)
			req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusCreated))
		})

		It("should login successfully", func() {
			payload := map[string]string{"username": "testuser", "password": "testpass"}
			body, _ := json.Marshal(payload)
			req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(httptest.NewRecorder(), req)

			req, _ = http.NewRequest("POST", "/users/login", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusOK))
			var response map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &response)
			Expect(response["token"]).NotTo(BeEmpty())
		})

		It("should fail login with wrong password", func() {
			payload := map[string]string{"username": "testuser", "password": "testpass"}
			body, _ := json.Marshal(payload)
			req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(httptest.NewRecorder(), req)

			wrongPayload := map[string]string{"username": "testuser", "password": "wrongpass"}
			body, _ = json.Marshal(wrongPayload)
			req, _ = http.NewRequest("POST", "/users/login", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusUnauthorized))
		})
	})

	Describe("Items", func() {
		It("should create an item", func() {
			payload := map[string]interface{}{"name": "Test Item", "description": "Test Description", "price": 19.99}
			body, _ := json.Marshal(payload)
			req, _ := http.NewRequest("POST", "/items", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusCreated))
		})

		It("should get all items", func() {
			payload := map[string]interface{}{"name": "Test Item", "description": "Test Description", "price": 19.99}
			body, _ := json.Marshal(payload)
			req, _ := http.NewRequest("POST", "/items", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(httptest.NewRecorder(), req)

			req, _ = http.NewRequest("GET", "/items", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusOK))
			var items []models.Item
			json.Unmarshal(w.Body.Bytes(), &items)
			Expect(len(items)).To(BeNumerically(">", 0))
		})
	})

	Describe("Cart Operations", func() {
		var token string

		BeforeEach(func() {
			payload := map[string]string{"username": "cartuser", "password": "cartpass"}
			body, _ := json.Marshal(payload)
			req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(httptest.NewRecorder(), req)

			req, _ = http.NewRequest("POST", "/users/login", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			var response map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &response)
			token = response["token"].(string)
		})

		It("should add item to cart", func() {
			itemPayload := map[string]interface{}{"name": "Cart Item", "description": "Test", "price": 29.99}
			body, _ := json.Marshal(itemPayload)
			req, _ := http.NewRequest("POST", "/items", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			var item models.Item
			json.Unmarshal(w.Body.Bytes(), &item)

			cartPayload := map[string]interface{}{"item_id": item.ID}
			body, _ = json.Marshal(cartPayload)
			req, _ = http.NewRequest("POST", "/carts", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+token)
			w = httptest.NewRecorder()
			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusCreated))
		})

		It("should get cart", func() {
			req, _ := http.NewRequest("GET", "/carts", nil)
			req.Header.Set("Authorization", "Bearer "+token)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusOK))
		})
	})

	Describe("Order Operations", func() {
		var token string

		BeforeEach(func() {
			payload := map[string]string{"username": "orderuser", "password": "orderpass"}
			body, _ := json.Marshal(payload)
			req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(httptest.NewRecorder(), req)

			req, _ = http.NewRequest("POST", "/users/login", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			var response map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &response)
			token = response["token"].(string)
		})

		It("should create order from cart", func() {
			itemPayload := map[string]interface{}{"name": "Order Item", "description": "Test", "price": 39.99}
			body, _ := json.Marshal(itemPayload)
			req, _ := http.NewRequest("POST", "/items", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			var item models.Item
			json.Unmarshal(w.Body.Bytes(), &item)

			cartPayload := map[string]interface{}{"item_id": item.ID}
			body, _ = json.Marshal(cartPayload)
			req, _ = http.NewRequest("POST", "/carts", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+token)
			router.ServeHTTP(httptest.NewRecorder(), req)

			req, _ = http.NewRequest("POST", "/orders", nil)
			req.Header.Set("Authorization", "Bearer "+token)
			w = httptest.NewRecorder()
			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusCreated))
		})

		It("should get all orders", func() {
			req, _ := http.NewRequest("GET", "/orders", nil)
			req.Header.Set("Authorization", "Bearer "+token)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusOK))
		})
	})
})
