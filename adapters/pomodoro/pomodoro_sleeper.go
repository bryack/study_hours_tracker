package pomodoro

import (
	"time"
)

type SleeperFunc func(duration time.Duration)

func (s SleeperFunc) Wait(duration time.Duration) {
	s(duration)
}

func PomodoroSleeper(duration time.Duration) {
	time.Sleep(duration)
}
