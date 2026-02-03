package main

import (
	"log"
	"os"

	"github.com/bryack/study_hours_tracker/adapters/cli"
	"github.com/bryack/study_hours_tracker/adapters/database"
	"github.com/bryack/study_hours_tracker/adapters/pomodoro"
	domainPomodoro "github.com/bryack/study_hours_tracker/domain/pomodoro"
)

func main() {
	store, err := database.SetupPostgres()
	if err != nil {
		log.Fatal(err)
	}

	pomodoroRunner := domainPomodoro.NewPomodoro(pomodoro.SleeperFunc(pomodoro.PomodoroSleeper))
	session := cli.NewStudySession(store, pomodoroRunner)
	tracker := cli.NewCLI(os.Stdin, os.Stdout, session)
	if err := tracker.Run(); err != nil {
		log.Fatal(err)
	}
}
