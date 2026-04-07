package models

type ReadingTracker struct {
	TrackerID  int     `json:"tracker_id"`
	BookID     *int    `json:"book_id"`
	StartDate  *string `json:"start_date"`
	FinishDate *string `json:"finish_date"`
	Status     *string `json:"status"`
	Rating     *int    `json:"rating"`
}
