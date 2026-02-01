package pomodoro

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type SpySleeper struct {
	WaitCalled int
	Duration   time.Duration
}

func (s *SpySleeper) Wait(duration time.Duration) {
	s.WaitCalled++
	s.Duration = duration
}

func TestPomodoro(t *testing.T) {

	t.Run("pomodoro waits on start up", func(t *testing.T) {
		sleeper := &SpySleeper{}
		p := NewPomodoro(sleeper)
		p.Wait(0)
		assert.Equal(t, 1, sleeper.WaitCalled)
	})
	t.Run("pomodoro schedules alert at correct time", func(t *testing.T) {
		sleeper := &SpySleeper{}
		duration := 25 * time.Minute
		p := NewPomodoro(sleeper)
		p.Wait(duration)
		assert.Equal(t, duration, sleeper.Duration)
	})
}
