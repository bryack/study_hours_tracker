package main

import (
	"log"
	"os"

	"github.com/bryack/study_hours_tracker/adapters/cli"
	"github.com/bryack/study_hours_tracker/adapters/database"
	"github.com/bryack/study_hours_tracker/adapters/pomodoro"
)

func main() {
	store, err := database.SetupPostgres()
	if err != nil {
		log.Fatal(err)
	}

	sleeper := pomodoro.SleeperFunc(pomodoro.PomodoroSleeper)

	tracker := cli.NewCLI(store, os.Stdin, os.Stdout, sleeper)
	tracker.Run()
}
