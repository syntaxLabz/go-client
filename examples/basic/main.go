package main

import (
	"fmt"
	"log"

	"github.com/yourorg/httpclient"
)

// User represents a user from the API
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func main() {
	fmt.Println("=== Ultra-Simple HTTP Client Examples ===\n")

	// Example 1: One-liner GET request
	fmt.Println("1. Simple GET request:")
	data, err := httpclient.GET("https://jsonplaceholder.typicode.com/users/1")
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		fmt.Printf("Response: %s\n\n", data)
	}

	// Example 2: JSON parsing made easy
	fmt.Println("2. GET with automatic JSON parsing:")
	var user User
	err = httpclient.JSON("GET", "https://jsonplaceholder.typicode.com/users/1", nil, &user)
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		fmt.Printf("User: %+v\n\n", user)
	}

	// Example 3: POST with JSON body
	fmt.Println("3. POST with JSON body:")
	newUser := User{Name: "John Doe", Email: "john@example.com"}
	var createdUser User
	err = httpclient.JSON("POST", "https://jsonplaceholder.typicode.com/users", newUser, &createdUser)
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		fmt.Printf("Created user: %+v\n\n", createdUser)
	}

	// Example 4: All HTTP methods
	fmt.Println("4. All HTTP methods:")
	
	// PUT
	updatedUser := User{ID: 1, Name: "Jane Doe", Email: "jane@example.com"}
	_, err = httpclient.PUT("https://jsonplaceholder.typicode.com/users/1", updatedUser)
	if err != nil {
		log.Printf("PUT Error: %v", err)
	} else {
		fmt.Println("PUT request successful")
	}

	// PATCH
	patchData := map[string]string{"name": "Updated Name"}
	_, err = httpclient.PATCH("https://jsonplaceholder.typicode.com/users/1", patchData)
	if err != nil {
		log.Printf("PATCH Error: %v", err)
	} else {
		fmt.Println("PATCH request successful")
	}

	// DELETE
	_, err = httpclient.DELETE("https://jsonplaceholder.typicode.com/users/1")
	if err != nil {
		log.Printf("DELETE Error: %v", err)
	} else {
		fmt.Println("DELETE request successful")
	}

	// HEAD
	err = httpclient.HEAD("https://jsonplaceholder.typicode.com/users/1")
	if err != nil {
		log.Printf("HEAD Error: %v", err)
	} else {
		fmt.Println("HEAD request successful")
	}

	fmt.Println()
}