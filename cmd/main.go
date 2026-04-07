package main

import (
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"

	"github.com/amber-bobby/Lab6_CMPS2242/db"
	"github.com/amber-bobby/Lab6_CMPS2242/handlers"
)

func main() {
	// Connect to the database
	db.Connect()

	// Register routes
	http.HandleFunc("/genres", handlers.GenreHandler)
	http.HandleFunc("/genres/", handlers.GenreHandler)

	http.HandleFunc("/books", handlers.BookHandler)
	http.HandleFunc("/books/", handlers.BookHandler)

	http.HandleFunc("/trackers", handlers.TrackerHandler)
	http.HandleFunc("/trackers/", handlers.TrackerHandler)

	port := ":8080"
	fmt.Printf("Server running on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
