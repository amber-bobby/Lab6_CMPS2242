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

// Router: /trackers and /trackers/{id}
func TrackerHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	hasID := len(parts) == 2 && parts[1] != ""

	if hasID {
		id, err := strconv.Atoi(parts[1])
		if err != nil {
			jsonError(w, "Invalid tracker ID", http.StatusBadRequest)
			return
		}
		switch r.Method {
		case http.MethodGet:
			getTrackerByID(w, id)
		case http.MethodPut:
			updateTracker(w, r, id)
		case http.MethodDelete:
			deleteTracker(w, id)
		default:
			jsonError(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
		return
	}

	switch r.Method {
	case http.MethodGet:
		getAllTrackers(w)
	case http.MethodPost:
		createTracker(w, r)
	default:
		jsonError(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getAllTrackers(w http.ResponseWriter) {
	rows, err := db.DB.Query(
		"SELECT tracker_id, book_id, start_date, finish_date, status, rating FROM reading_tracker",
	)
	if err != nil {
		jsonError(w, "Failed to fetch trackers", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	trackers := []models.ReadingTracker{}
	for rows.Next() {
		var t models.ReadingTracker
		if err := rows.Scan(&t.TrackerID, &t.BookID, &t.StartDate, &t.FinishDate, &t.Status, &t.Rating); err != nil {
			jsonError(w, "Failed to scan tracker", http.StatusInternalServerError)
			return
		}
		trackers = append(trackers, t)
	}
	jsonResponse(w, trackers, http.StatusOK)
}

func getTrackerByID(w http.ResponseWriter, id int) {
	var t models.ReadingTracker
	err := db.DB.QueryRow(
		"SELECT tracker_id, book_id, start_date, finish_date, status, rating FROM reading_tracker WHERE tracker_id = $1", id,
	).Scan(&t.TrackerID, &t.BookID, &t.StartDate, &t.FinishDate, &t.Status, &t.Rating)
	if err == sql.ErrNoRows {
		jsonError(w, "Tracker not found", http.StatusNotFound)
		return
	} else if err != nil {
		jsonError(w, "Failed to fetch tracker", http.StatusInternalServerError)
		return
	}
	jsonResponse(w, t, http.StatusOK)
}

func createTracker(w http.ResponseWriter, r *http.Request) {
	var t models.ReadingTracker
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		jsonError(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	err := db.DB.QueryRow(
		"INSERT INTO reading_tracker (book_id, start_date, finish_date, status, rating) VALUES ($1, $2, $3, $4, $5) RETURNING tracker_id",
		t.BookID, t.StartDate, t.FinishDate, t.Status, t.Rating,
	).Scan(&t.TrackerID)
	if err != nil {
		jsonError(w, "Failed to create tracker", http.StatusInternalServerError)
		return
	}
	jsonResponse(w, t, http.StatusCreated)
}

func updateTracker(w http.ResponseWriter, r *http.Request, id int) {
	var t models.ReadingTracker
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		jsonError(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	result, err := db.DB.Exec(
		"UPDATE reading_tracker SET book_id = $1, start_date = $2, finish_date = $3, status = $4, rating = $5 WHERE tracker_id = $6",
		t.BookID, t.StartDate, t.FinishDate, t.Status, t.Rating, id,
	)
	if err != nil {
		jsonError(w, "Failed to update tracker", http.StatusInternalServerError)
		return
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		jsonError(w, "Tracker not found", http.StatusNotFound)
		return
	}
	t.TrackerID = id
	jsonResponse(w, t, http.StatusOK)
}

func deleteTracker(w http.ResponseWriter, id int) {
	result, err := db.DB.Exec("DELETE FROM reading_tracker WHERE tracker_id = $1", id)
	if err != nil {
		jsonError(w, "Failed to delete tracker", http.StatusInternalServerError)
		return
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		jsonError(w, "Tracker not found", http.StatusNotFound)
		return
	}
	jsonResponse(w, map[string]string{"message": "Tracker deleted successfully"}, http.StatusOK)
}
