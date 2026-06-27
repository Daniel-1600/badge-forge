package main

import (
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"tahrir-go/internal/db"
	"tahrir-go/internal/handlers"
	"tahrir-go/internal/rules"
	"tahrir-go/internal/rules/workers"
)

// @title           Tahrir API
// @version         1.0
// @description     Badge assertion API
// @host            localhost:8080
// @BasePath        /

func main() {
	_ = godotenv.Load()
	dsn := "host=" + os.Getenv("DB_HOST") + " user=" + os.Getenv("DB_USER") + " password=" + os.Getenv("DB_PASSWORD") + " dbname=" + os.Getenv("DB_NAME") + " port=" + os.Getenv("DB_PORT") + " sslmode=disable"
	conn, err := db.Connect(dsn)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	log.Println("connected to Tahrir database")

	eventChannel := make(chan rules.Event)

	w := worker.Worker{
		Events: eventChannel,
		DB:     conn,
		Rules: []rules.Rule{
			&rules.MilestoneRule{Threshold: 3, DB: conn},
		},
	}

	w.Start()

	// set up routes
	http.HandleFunc("GET /persons", handlers.GetPersons(conn))
	http.HandleFunc("GET /persons/{nickname}", handlers.GetPersonByNickname(conn))
	http.HandleFunc("GET /persons/id/{id}", handlers.GetPersonByID(conn))
	http.HandleFunc("GET /badges", handlers.GetBadges(conn))
	http.HandleFunc("GET /badges/{id}", handlers.GetBadgeByID(conn))
	http.HandleFunc("GET /assertions/{id}", handlers.GetAssertionByID(conn))
	http.HandleFunc("GET /persons/nickname/{person_nickname}/badges", handlers.GetAssertionsByPersonNickname(conn))
	http.HandleFunc("POST /badges", handlers.CreateBadge(conn))
	http.HandleFunc("POST /assertions", handlers.PostAssertion(conn, eventChannel))

	//start the server
	log.Println("starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("server failed: %v", err)
	}

}
