package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
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

var db *sql.DB

// In-memory store for a single user (for simplicity)
var username = "admin"
var password = "password"

// Initialize and connect to the SQLite database
func init() {
	var err error
	db, err = initDB()
	if err != nil {
		log.Fatalf("Could not initialize database: %v", err)
	}
}

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
	books, err := getBooksFromDB(db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

// Get a single book by ID
func getBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	book, err := getBookFromDB(db, params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if book == nil {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

// Create a new book
func createBook(w http.ResponseWriter, r *http.Request) {
	var book Book
	_ = json.NewDecoder(r.Body).Decode(&book)
	book.ID = uuid.New().String() // Generate a new unique ID for the book

	err := createBookInDB(db, book)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

// Delete a book by ID
func deleteBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	err := deleteBookFromDB(db, params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"result": "success"})
}

func main() {
	// Initialize router
	router := mux.NewRouter()

	// Login endpoint to generate JWT token
	router.HandleFunc("/login", login).Methods("POST")

	// Secured routes with JWT authentication
	apiRouter := router.PathPrefix("/api").Subrouter()
	apiRouter.Use(jwtAuthMiddleware)

	// Route handlers & endpoints
	apiRouter.HandleFunc("/books", getBooks).Methods("GET")
	apiRouter.HandleFunc("/books/{id}", getBook).Methods("GET")
	apiRouter.HandleFunc("/books", createBook).Methods("POST")
	apiRouter.HandleFunc("/books/{id}", deleteBook).Methods("DELETE")

	// Start server
	log.Println("Server running on port 10000")
	log.Fatal(http.ListenAndServe(":10000", router))
}
