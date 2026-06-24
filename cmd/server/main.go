package main

import (
	"log"
	"net/http"
	"tahrir-go/internal/db"
	"tahrir-go/internal/handlers"
)

// @title           Tahrir API
// @version         1.0
// @description     Badge assertion API
// @host            localhost:8080
// @BasePath        /

func main() {

	dsn := "host=localhost user=tahrir password=tahrir dbname=tahrir port=5432 sslmode=disable"
	conn, err := db.Connect(dsn)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	log.Println("connected to Tahrir database")

	// 3. set up routes
	http.HandleFunc("GET /persons", handlers.GetPersons(conn))
	http.HandleFunc("GET /persons/{nickname}", handlers.GetPersonByNickname(conn))
	http.HandleFunc("GET /persons/id/{id}", handlers.GetPersonByID(conn))
	http.HandleFunc("GET /badges", handlers.GetBadges(conn))
	http.HandleFunc("GET /badges/{id}", handlers.GetBadgeByID(conn))
	http.HandleFunc("GET /assertions/{id}", handlers.GetAssertionsByID(conn))
	http.HandleFunc("GET /persons/nickname/{person_nickname}/badges", handlers.GetAssertionsByPersonNickname(conn))
	http.HandleFunc("POST /badges", handlers.CreateBadge(conn))
	http.HandleFunc("POST /assertions", handlers.PostAssertion(conn))

	// 4. start the server
	log.Println("starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("server failed: %v", err)
	}

}
