package main

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

// Setup a helper function to create a test database connection
func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:") // In-memory database for testing
	if err != nil {
		t.Fatalf("Could not create test database: %v", err)
	}

	// Create the books table in the in-memory database
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS books (
		id TEXT PRIMARY KEY,
		title TEXT,
		author TEXT,
		year TEXT
	);`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		t.Fatalf("Could not create books table: %v", err)
	}

	return db
}

func TestCreateBookInDB(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	book := Book{
		ID:     "1",
		Title:  "Test Book",
		Author: "Test Author",
		Year:   "2023",
	}

	err := createBookInDB(db, book)
	assert.NoError(t, err, "Expected no error when inserting a book")

	// Verify the book was inserted
	var title string
	err = db.QueryRow("SELECT title FROM books WHERE id = ?", book.ID).Scan(&title)
	assert.NoError(t, err, "Expected no error when fetching the inserted book")
	assert.Equal(t, book.Title, title, "Expected inserted book title to match")
}

func TestGetBooksFromDB(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	book := Book{
		ID:     "1",
		Title:  "Test Book",
		Author: "Test Author",
		Year:   "2023",
	}

	err := createBookInDB(db, book)
	assert.NoError(t, err, "Expected no error when inserting a book")

	books, err := getBooksFromDB(db)
	assert.NoError(t, err, "Expected no error when getting books")
	assert.Len(t, books, 1, "Expected 1 book in database")
	assert.Equal(t, book.Title, books[0].Title, "Expected title to be what we inserted")
}

func TestGetBookFromDB(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Insert a test book
	book := Book{
		ID:     "1",
		Title:  "Test Book",
		Author: "Test Author",
		Year:   "2023",
	}
	err := createBookInDB(db, book)
	assert.NoError(t, err, "Expected no error when inserting a book")

	// Fetch the book by ID
	fetchedBook, err := getBookFromDB(db, book.ID)
	assert.NoError(t, err, "Expected no error when fetching the book")
	assert.NotNil(t, fetchedBook, "Expected to find the book in the database")
	assert.Equal(t, book.Title, fetchedBook.Title, "Expected fetched book title to match")
}

func TestDeleteBookFromDB(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Insert a test book
	book := Book{
		ID:     "1",
		Title:  "Test Book",
		Author: "Test Author",
		Year:   "2023",
	}
	err := createBookInDB(db, book)
	assert.NoError(t, err, "Expected no error when inserting a book")

	// Delete the book
	err = deleteBookFromDB(db, book.ID)
	assert.NoError(t, err, "Expected no error when deleting the book")

	// Verify the book was deleted
	fetchedBook, err := getBookFromDB(db, book.ID)
	assert.NoError(t, err, "Expected no error when fetching the deleted book")
	assert.Nil(t, fetchedBook, "Expected the book to be deleted")
}
