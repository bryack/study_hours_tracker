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
	t.Run("pomodoro waits correct time", func(t *testing.T) {
		sleeper := &SpySleeper{}
		p := NewPomodoro(sleeper)
		p.Start()
		assert.Equal(t, 1, sleeper.WaitCalled)
		assert.Equal(t, p.duration, sleeper.Duration)
	})
}
