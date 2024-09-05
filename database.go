package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

// Initialize SQLite database connection
func initDB() (*sql.DB, error) {
	log.Println("Initializing DB..")
	db, err := sql.Open("sqlite3", "./books.db")
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
