package pomodoro

type Sleeper interface {
	Wait()
}

type Pomodoro struct {
	sleeper Sleeper
}

func NewPomodoro(sleeper Sleeper) *Pomodoro {
	return &Pomodoro{
		sleeper: sleeper,
	}
}

func (p *Pomodoro) Start() {
	p.sleeper.Wait()
}
