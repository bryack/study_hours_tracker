package main

import (
	"log"
	"net/http"
	"os"

	"github.com/bryack/study_hours_tracker/adapters/database"
	"github.com/bryack/study_hours_tracker/adapters/server"
)

func main() {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		connStr = "postgres://postgres:pass@localhost:5432/mydb?sslmode=disable"
	}
	store, err := database.NewPostgresSubjectStore(connStr)
	if err != nil {
		log.Fatalf("could not connect to database: %v", err)
	}
	svr := server.NewStudyServer(store)

	log.Fatal(http.ListenAndServe(":5000", svr))
}
