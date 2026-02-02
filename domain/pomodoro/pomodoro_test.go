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

func (s *SpySleeper) Wait() {
	s.WaitCalled++
	s.Duration = 25 * time.Minute
}

func TestPomodoro(t *testing.T) {
	t.Run("pomodoro waits correct time", func(t *testing.T) {
		sleeper := &SpySleeper{}
		p := NewPomodoro(sleeper)
		p.Start()
		assert.Equal(t, 1, sleeper.WaitCalled)
		assert.Equal(t, 25*time.Minute, sleeper.Duration)
	})
}
