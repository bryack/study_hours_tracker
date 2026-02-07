package main

import (
	"log"
	"net/http"

	"github.com/bryack/study_hours_tracker/adapters/database"
	"github.com/bryack/study_hours_tracker/adapters/server"
)

const defaultPort = ":5000"

func main() {
	store, err := database.SetupPostgres()
	if err != nil {
		log.Fatal(err)
	}

	svr, err := server.NewStudyServer(store)
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(http.ListenAndServe(defaultPort, svr))
}
