package main

import (
	"fmt"
	"log"
	"os"

	"github.com/bryack/study_hours_tracker/adapters/cli"
	"github.com/bryack/study_hours_tracker/adapters/database"
)

func main() {
	fmt.Println("Let's study")
	fmt.Println("Type {subject} {hours} to track hours")

	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		connStr = "postgres://postgres:pass@localhost:5432/mydb?sslmode=disable"
	}
	store, err := database.NewPostgresSubjectStore(connStr)
	if err != nil {
		log.Fatalf("could not connect to database: %v", err)
	}

	tracker := cli.NewCLI(store, os.Stdin)
	tracker.Run()
}
