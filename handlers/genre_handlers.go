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

// Router: /genres and /genres/{id}
func GenreHandler(w http.ResponseWriter, r *http.Request) {
	// Check if an ID is in the path e.g. /genres/1
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	hasID := len(parts) == 2 && parts[1] != ""

	if hasID {
		id, err := strconv.Atoi(parts[1])
		if err != nil {
			jsonError(w, "Invalid genre ID", http.StatusBadRequest)
			return
		}
		switch r.Method {
		case http.MethodGet:
			getGenreByID(w, id)
		case http.MethodPut:
			updateGenre(w, r, id)
		case http.MethodDelete:
			deleteGenre(w, id)
		default:
			jsonError(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
		return
	}

	switch r.Method {
	case http.MethodGet:
		getAllGenres(w)
	case http.MethodPost:
		createGenre(w, r)
	default:
		jsonError(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getAllGenres(w http.ResponseWriter) {
	rows, err := db.DB.Query("SELECT genre_id, name FROM genres")
	if err != nil {
		jsonError(w, "Failed to fetch genres", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	genres := []models.Genre{}
	for rows.Next() {
		var g models.Genre
		if err := rows.Scan(&g.GenreID, &g.Name); err != nil {
			jsonError(w, "Failed to scan genre", http.StatusInternalServerError)
			return
		}
		genres = append(genres, g)
	}
	jsonResponse(w, genres, http.StatusOK)
}

func getGenreByID(w http.ResponseWriter, id int) {
	var g models.Genre
	err := db.DB.QueryRow("SELECT genre_id, name FROM genres WHERE genre_id = $1", id).
		Scan(&g.GenreID, &g.Name)
	if err == sql.ErrNoRows {
		jsonError(w, "Genre not found", http.StatusNotFound)
		return
	} else if err != nil {
		jsonError(w, "Failed to fetch genre", http.StatusInternalServerError)
		return
	}
	jsonResponse(w, g, http.StatusOK)
}

func createGenre(w http.ResponseWriter, r *http.Request) {
	var g models.Genre
	if err := json.NewDecoder(r.Body).Decode(&g); err != nil {
		jsonError(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	err := db.DB.QueryRow(
		"INSERT INTO genres (name) VALUES ($1) RETURNING genre_id",
		g.Name,
	).Scan(&g.GenreID)
	if err != nil {
		jsonError(w, "Failed to create genre", http.StatusInternalServerError)
		return
	}
	jsonResponse(w, g, http.StatusCreated)
}

func updateGenre(w http.ResponseWriter, r *http.Request, id int) {
	var g models.Genre
	if err := json.NewDecoder(r.Body).Decode(&g); err != nil {
		jsonError(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	result, err := db.DB.Exec("UPDATE genres SET name = $1 WHERE genre_id = $2", g.Name, id)
	if err != nil {
		jsonError(w, "Failed to update genre", http.StatusInternalServerError)
		return
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		jsonError(w, "Genre not found", http.StatusNotFound)
		return
	}
	g.GenreID = id
	jsonResponse(w, g, http.StatusOK)
}

func deleteGenre(w http.ResponseWriter, id int) {
	result, err := db.DB.Exec("DELETE FROM genres WHERE genre_id = $1", id)
	if err != nil {
		jsonError(w, "Failed to delete genre", http.StatusInternalServerError)
		return
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		jsonError(w, "Genre not found", http.StatusNotFound)
		return
	}
	jsonResponse(w, map[string]string{"message": "Genre deleted successfully"}, http.StatusOK)
}
