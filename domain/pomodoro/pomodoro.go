package pomodoro

import "time"

type Sleeper interface {
	Wait(duration time.Duration)
}

type Pomodoro struct {
	sleeper Sleeper
}

func NewPomodoro(sleeper Sleeper) *Pomodoro {
	return &Pomodoro{
		sleeper: sleeper,
	}
}

func (p *Pomodoro) Wait(duration time.Duration) {
	p.sleeper.Wait(duration)
}
