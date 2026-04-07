package models

type Book struct {
	BookID        int    `json:"book_id"`
	Title         string `json:"title"`
	Author        string `json:"author"`
	GenreID       *int   `json:"genre_id"`
	PublishedYear *int   `json:"published_year"`
	Rating        string `json:"rating"`
}
