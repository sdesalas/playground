package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

// Simple struct to represent a "book"
type Book struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
	Year   string `json:"year"`
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// In-memory data store (for simplicity)
var books []Book

// In-memory store for a single user (for simplicity)
var username = "admin"
var password = "password"

// JWT login handler
func login(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Check if the provided credentials are correct
	if creds.Username != username || creds.Password != password {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Generate JWT token
	tokenString, err := GenerateJWTToken(creds.Username)
	if err != nil {
		http.Error(w, "Could not generate token", http.StatusInternalServerError)
		return
	}

	// Send the token as a response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}

// Middleware to validate JWT token
func jwtAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Middleware jwtAuthMiddleware..")

		authHeader := r.Header.Get("Authorization")

		// Check if Authorization header is present and formatted correctly
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Missing or invalid Authorization header", http.StatusUnauthorized)
			return
		}

		// Extract token from the header
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Validate the JWT token
		_, err := ValidateJWTToken(tokenString)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		// If token is valid, proceed to the next handler
		next.ServeHTTP(w, r)
	})
}

// Get all books
func getBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

// Get a single book by ID
func getBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r) // Get URL params
	for _, item := range books {
		if item.ID == params["id"] {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	http.Error(w, "Book not found", http.StatusNotFound)
}

// Create a new book
func createBook(w http.ResponseWriter, r *http.Request) {
	var book Book
	_ = json.NewDecoder(r.Body).Decode(&book)
	book.ID = fmt.Sprintf("%d", len(books)+1) // Simple ID generation
	books = append(books, book)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

// Delete a book by ID
func deleteBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for index, item := range books {
		if item.ID == params["id"] {
			books = append(books[:index], books[index+1:]...) // Remove the book
			break
		}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}


func main() {
	// Initialize router
	router := mux.NewRouter()

	// Sample books
	books = append(books, Book{ID: "1", Title: "1984", Author: "George Orwell", Year: "1949"})
	books = append(books, Book{ID: "2", Title: "The Great Gatsby", Author: "F. Scott Fitzgerald", Year: "1925"})

	// Login endpoint to generate JWT token
	router.HandleFunc("/login", login).Methods("POST")

	// Secured routes with JWT authentication
	apiRouter := router.PathPrefix("/api").Subrouter()
	apiRouter.Use(jwtAuthMiddleware)

	// Route handlers & endpoints
	router.HandleFunc("/api/books", getBooks).Methods("GET")
	router.HandleFunc("/api/books/{id}", getBook).Methods("GET")
	router.HandleFunc("/api/books", createBook).Methods("POST")
	router.HandleFunc("/api/books/{id}", deleteBook).Methods("DELETE")

	// Start server
	log.Println("Server running on port 8000")
	log.Fatal(http.ListenAndServe(":8000", router))
}