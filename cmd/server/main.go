package main

import (
	"log"
	"net/http"
	"tahrir-go/internal/db"
	"tahrir-go/internal/handlers"
)

func main() {

	dsn := "host=localhost user=tahrir password=tahrir dbname=tahrir port=5432 sslmode=disable"
	conn, err := db.Connect(dsn)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	log.Println("connected to Tahrir database")

	// 3. set up routes
	http.HandleFunc("/persons", handlers.GetPersons(conn))

	// 4. start the server
	log.Println("starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("server failed: %v", err)
	}

}
