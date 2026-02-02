package pomodoro

import "time"

type Sleeper interface {
	Wait(duration time.Duration)
}

type Pomodoro struct {
	sleeper  Sleeper
	duration time.Duration
}

func NewPomodoro(sleeper Sleeper) *Pomodoro {
	return &Pomodoro{
		sleeper:  sleeper,
		duration: 25 * time.Minute,
	}
}

func (p *Pomodoro) Start() {
	p.sleeper.Wait(p.duration)
}
