// Package pomodoro provides infrastructure implementations for timer functionality.
package pomodoro

import (
	"time"
)

// SleeperFunc is a function type that implements the Sleeper interface.
// It allows any function with the signature func(time.Duration) to be used as a Sleeper.
//
// Example usage:
//
//	sleeper := SleeperFunc(PomodoroSleeper)
//	pomodoro := domain.NewPomodoro(sleeper)
type SleeperFunc func(duration time.Duration)

// Wait calls the underlying function with the given duration.
func (s SleeperFunc) Wait(duration time.Duration) {
	s(duration)
}

// PomodoroSleeper is a Sleeper implementation that uses time.Sleep.
// It blocks for the specified duration.
func PomodoroSleeper(duration time.Duration) {
	time.Sleep(duration)
}
