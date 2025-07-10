package main

import (
	"fmt"
	"log"
	"time"

	"github.com/yourorg/httpclient"
)

// Service clients for microservice architecture
type UserService struct {
	client httpclient.Client
}

type OrderService struct {
	client httpclient.Client
}

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Order struct {
	ID     int `json:"id"`
	UserID int `json:"userId"`
	Total  int `json:"total"`
}

func NewUserService(baseURL string) *UserService {
	client := httpclient.New().
		WithBaseURL(baseURL).
		WithTimeout(5 * time.Second).
		WithRetries(3).
		WithAuth("user-service-token").
		WithHeader("X-Service", "user-service").
		WithMetrics(true)

	return &UserService{client: client}
}

func NewOrderService(baseURL string) *OrderService {
	client := httpclient.New().
		WithBaseURL(baseURL).
		WithTimeout(10 * time.Second).
		WithRetries(5).
		WithAuth("order-service-token").
		WithHeader("X-Service", "order-service").
		WithCircuitBreaker(5, 60*time.Second).
		WithMetrics(true)

	return &OrderService{client: client}
}

func (us *UserService) GetUser(id int) (*User, error) {
	var user User
	err := us.client.JSON("GET", fmt.Sprintf("/users/%d", id), nil, &user)
	if err != nil {
		return nil, fmt.Errorf("failed to get user %d: %w", id, err)
	}
	return &user, nil
}

func (us *UserService) CreateUser(user *User) (*User, error) {
	var created User
	err := us.client.JSON("POST", "/users", user, &created)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	return &created, nil
}

func (os *OrderService) GetOrder(id int) (*Order, error) {
	var order Order
	err := os.client.JSON("GET", fmt.Sprintf("/orders/%d", id), nil, &order)
	if err != nil {
		return nil, fmt.Errorf("failed to get order %d: %w", id, err)
	}
	return &order, nil
}

func (os *OrderService) CreateOrder(order *Order) (*Order, error) {
	var created Order
	err := os.client.JSON("POST", "/orders", order, &created)
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}
	return &created, nil
}

func main() {
	fmt.Println("=== Microservice Communication Example ===\n")

	// Initialize service clients
	userService := NewUserService("https://jsonplaceholder.typicode.com")
	orderService := NewOrderService("https://jsonplaceholder.typicode.com")

	// Use the services
	fmt.Println("1. Getting user from user service:")
	user, err := userService.GetUser(1)
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		fmt.Printf("User: %+v\n\n", user)
	}

	fmt.Println("2. Creating new user:")
	newUser := &User{Name: "Service User", Email: "service@example.com"}
	created, err := userService.CreateUser(newUser)
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		fmt.Printf("Created user: %+v\n\n", created)
	}

	// Simulate order service usage
	fmt.Println("3. Order service configured with circuit breaker and extended timeout")
	fmt.Println("   (Would work with real order service endpoints)\n")

	// Example of service-to-service communication
	fmt.Println("4. Service-to-service communication pattern:")
	fmt.Println("   - User service: 5s timeout, 3 retries, metrics enabled")
	fmt.Println("   - Order service: 10s timeout, 5 retries, circuit breaker, metrics enabled")
	fmt.Println("   - Both services have authentication and service identification headers")
}