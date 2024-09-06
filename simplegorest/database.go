package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/lib/pq"
)

// Initialize SQLite database connection
func initDB() (*sql.DB, error) {
	log.Println("Initializing DB..")
	connStr := os.Getenv("DB_CONNECTIONSTRING")
	log.Printf("Connecting to %s...", strings.Split(connStr, "@")[0])
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	// Create the books table if it doesn't exist
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS books (
		id TEXT PRIMARY KEY,
		title TEXT,
		author TEXT,
		year TEXT
	);`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		return nil, fmt.Errorf("error creating table: %v", err)
	}

	rows, err := db.Query("SELECT * FROM books LIMIT 5")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var id, title, author string
	var year int
	for rows.Next() {
		err := rows.Scan(&id, &title, &author, &year)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("id: %s, title: %s, author: %s, year: %d", id, title, author, year)
	}

	return db, nil
}

// Insert a new book into the database
func createBookInDB(db *sql.DB, book Book) error {
	insertBookQuery := `INSERT INTO books (id, title, author, year) VALUES (?, ?, ?, ?)`
	_, err := db.Exec(insertBookQuery, book.ID, book.Title, book.Author, book.Year)
	if err != nil {
		return fmt.Errorf("error inserting book: %v", err)
	}
	return nil
}

// Get all books from the database
func getBooksFromDB(db *sql.DB) ([]Book, error) {
	rows, err := db.Query("SELECT id, title, author, year FROM books")
	if err != nil {
		return nil, fmt.Errorf("error fetching books: %v", err)
	}
	defer rows.Close()

	var books []Book
	for rows.Next() {
		var book Book
		err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Year)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}
		books = append(books, book)
	}
	return books, nil
}

// Get a single book by ID from the database
func getBookFromDB(db *sql.DB, id string) (*Book, error) {
	var book Book
	query := `SELECT id, title, author, year FROM books WHERE id = ?`
	err := db.QueryRow(query, id).Scan(&book.ID, &book.Title, &book.Author, &book.Year)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error fetching book by ID: %v", err)
	}
	return &book, nil
}

// Delete a book from the database
func deleteBookFromDB(db *sql.DB, id string) error {
	deleteQuery := `DELETE FROM books WHERE id = ?`
	_, err := db.Exec(deleteQuery, id)
	if err != nil {
		return fmt.Errorf("error deleting book: %v", err)
	}
	return nil
}
