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

	store, err := database.SetupPostgres()
	if err != nil {
		log.Fatal(err)
	}

	tracker := cli.NewCLI(store, os.Stdin)
	tracker.Run()
}
