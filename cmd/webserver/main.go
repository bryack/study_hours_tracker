package main

import (
	"log"
	"net/http"

	"github.com/bryack/study_hours_tracker/adapters/database"
	"github.com/bryack/study_hours_tracker/adapters/pomodoro"
	"github.com/bryack/study_hours_tracker/adapters/server"
	"github.com/bryack/study_hours_tracker/domain"
	domainPomodoro "github.com/bryack/study_hours_tracker/domain/pomodoro"
)

const defaultPort = ":5000"

func main() {
	store, err := database.SetupPostgres()
	if err != nil {
		log.Fatal(err)
	}

	alerter := pomodoro.Alerter{
		ScheduleFunc: pomodoro.RealScheduleAlert,
		WaitFunc:     pomodoro.RealWait,
	}

	pomodoroRunner := domainPomodoro.NewPomodoro(alerter)
	session := domain.NewStudySession(store, pomodoroRunner)

	svr, err := server.NewStudyServer(store, session)
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(http.ListenAndServe(defaultPort, svr))
}
