package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/amber-bobby/Lab6_CMPS2242/db"
	"github.com/amber-bobby/Lab6_CMPS2242/models"
)

// Router: /books and /books/{id}
func BookHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	hasID := len(parts) == 2 && parts[1] != ""

	if hasID {
		id, err := strconv.Atoi(parts[1])
		if err != nil {
			jsonError(w, "Invalid book ID", http.StatusBadRequest)
			return
		}
		switch r.Method {
		case http.MethodGet:
			getBookByID(w, id)
		case http.MethodPut:
			updateBook(w, r, id)
		case http.MethodDelete:
			deleteBook(w, id)
		default:
			jsonError(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
		return
	}

	switch r.Method {
	case http.MethodGet:
		getAllBooks(w)
	case http.MethodPost:
		createBook(w, r)
	default:
		jsonError(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getAllBooks(w http.ResponseWriter) {
	rows, err := db.DB.Query("SELECT book_id, title, author, genre_id, published_year, rating FROM books")
	if err != nil {
		jsonError(w, "Failed to fetch books", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	books := []models.Book{}
	for rows.Next() {
		var b models.Book
		if err := rows.Scan(&b.BookID, &b.Title, &b.Author, &b.GenreID, &b.PublishedYear, &b.Rating); err != nil {
			jsonError(w, "Failed to scan book", http.StatusInternalServerError)
			return
		}
		books = append(books, b)
	}
	jsonResponse(w, books, http.StatusOK)
}

func getBookByID(w http.ResponseWriter, id int) {
	var b models.Book
	err := db.DB.QueryRow(
		"SELECT book_id, title, author, genre_id, published_year, rating FROM books WHERE book_id = $1", id,
	).Scan(&b.BookID, &b.Title, &b.Author, &b.GenreID, &b.PublishedYear, &b.Rating)
	if err == sql.ErrNoRows {
		jsonError(w, "Book not found", http.StatusNotFound)
		return
	} else if err != nil {
		jsonError(w, "Failed to fetch book", http.StatusInternalServerError)
		return
	}
	jsonResponse(w, b, http.StatusOK)
}

func createBook(w http.ResponseWriter, r *http.Request) {
	var b models.Book
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		jsonError(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	err := db.DB.QueryRow(
		"INSERT INTO books (title, author, genre_id, published_year, rating) VALUES ($1, $2, $3, $4, $5) RETURNING book_id",
		b.Title, b.Author, b.GenreID, b.PublishedYear, b.Rating,
	).Scan(&b.BookID)
	if err != nil {
		jsonError(w, "Failed to create book", http.StatusInternalServerError)
		return
	}
	jsonResponse(w, b, http.StatusCreated)
}

func updateBook(w http.ResponseWriter, r *http.Request, id int) {
	var b models.Book
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		jsonError(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	result, err := db.DB.Exec(
		"UPDATE books SET title = $1, author = $2, genre_id = $3, published_year = $4, rating = $5 WHERE book_id = $6",
		b.Title, b.Author, b.GenreID, b.PublishedYear, b.Rating, id,
	)
	if err != nil {
		jsonError(w, "Failed to update book", http.StatusInternalServerError)
		return
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		jsonError(w, "Book not found", http.StatusNotFound)
		return
	}
	b.BookID = id
	jsonResponse(w, b, http.StatusOK)
}

func deleteBook(w http.ResponseWriter, id int) {
	result, err := db.DB.Exec("DELETE FROM books WHERE book_id = $1", id)
	if err != nil {
		jsonError(w, "Failed to delete book", http.StatusInternalServerError)
		return
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		jsonError(w, "Book not found", http.StatusNotFound)
		return
	}
	jsonResponse(w, map[string]string{"message": "Book deleted successfully"}, http.StatusOK)
}
