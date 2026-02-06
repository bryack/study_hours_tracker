// Package pomodoro provides timer functionality for focused study sessions.
package pomodoro

import (
	"io"
	"time"
)

const DefaultPomodoroDuration = 25 * time.Minute

// PomodoroAlerter represents a timer that can wait for a specified duration.
type PomodoroAlerter interface {
	ScheduleAlert(duration time.Duration, message string, out io.Writer)
	Wait(duration time.Duration)
}

// Pomodoro represents a timer for focused study sessions using the Pomodoro Technique.
type Pomodoro struct {
	alerter  PomodoroAlerter
	duration time.Duration
}

// NewPomodoro creates a new Pomodoro timer with a default duration of 25 minutes.
func NewPomodoro(alerter PomodoroAlerter) *Pomodoro {
	return &Pomodoro{
		alerter:  alerter,                 // timer implementation
		duration: DefaultPomodoroDuration, // length of pomodoro session
	}
}

// Start begins the Pomodoro timer and waits for the configured duration.
func (p *Pomodoro) Start(out io.Writer) {
	p.alerter.ScheduleAlert(0, "Session started. Stay focused!", out)
	p.alerter.ScheduleAlert(p.duration/2, "Halfway there! Keep it up.", out)
	p.alerter.ScheduleAlert(p.duration, "Time's up! Recording your hour...", out)
	p.alerter.Wait(p.duration)
}
