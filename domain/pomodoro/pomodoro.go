// Package pomodoro provides timer functionality for focused study sessions.
package pomodoro

import "time"

const DefaultPomodoroDuration = 25 * time.Minute

// Sleeper represents a timer that can wait for a specified duration.
type Sleeper interface {
	Wait(duration time.Duration)
}

// Pomodoro represents a timer for focused study sessions using the Pomodoro Technique.
type Pomodoro struct {
	sleeper  Sleeper
	duration time.Duration
}

// NewPomodoro creates a new Pomodoro timer with a default duration of 25 minutes.
func NewPomodoro(sleeper Sleeper) *Pomodoro {
	return &Pomodoro{
		sleeper:  sleeper,                 // timer implementation
		duration: DefaultPomodoroDuration, // length of pomodoro session
	}
}

// Start begins the Pomodoro timer and waits for the configured duration.
func (p *Pomodoro) Start() {
	p.sleeper.Wait(p.duration)
}
