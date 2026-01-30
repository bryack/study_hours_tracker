package main

import (
	"log"
	"net/http"

	"github.com/bryack/study_hours_tracker/adapters/database"
	"github.com/bryack/study_hours_tracker/adapters/server"
)

func main() {
	store, err := database.SetupPostgres()
	if err != nil {
		log.Fatal(err)
	}

	svr := server.NewStudyServer(store)

	log.Fatal(http.ListenAndServe(":5000", svr))
}
