package pomodoro

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type SpySleeper struct {
	WaitCalled   int
	LastDuration time.Duration
}

func (s *SpySleeper) Wait(duration time.Duration) {
	s.WaitCalled++
	s.LastDuration = duration
}

func TestPomodoro(t *testing.T) {
	tests := []struct {
		name             string
		startCalls       int
		expectedCalls    int
		expectedDuration time.Duration
	}{
		{
			name:             "single start waits correct time",
			startCalls:       1,
			expectedCalls:    1,
			expectedDuration: DefaultPomodoroDuration,
		},
		{
			name:             "multiple starts call wait multiple times",
			startCalls:       3,
			expectedCalls:    3,
			expectedDuration: DefaultPomodoroDuration,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			out := &bytes.Buffer{}
			sleeper := &SpySleeper{}
			p := NewPomodoro(sleeper)
			for range tt.startCalls {
				p.Start(out)
			}
			assert.Equal(t, tt.expectedCalls, sleeper.WaitCalled)
			assert.Equal(t, DefaultPomodoroDuration, sleeper.LastDuration)
		})
	}
}
