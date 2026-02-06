package pomodoro

import (
	"bytes"
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type ScheduledAlert struct {
	At      time.Duration
	Message string
}

type SpyScheduleAlerter struct {
	Alerts     []ScheduledAlert
	WaitCalled int
}

func (s *SpyScheduleAlerter) ScheduleAlert(duration time.Duration, message string, out io.Writer) {
	s.Alerts = append(s.Alerts, ScheduledAlert{At: duration, Message: message})
}

func (s *SpyScheduleAlerter) Wait(duration time.Duration) {
	s.WaitCalled++
}

func TestPomodoro_Start(t *testing.T) {
	testcases := []ScheduledAlert{
		{0, "Session started. Stay focused!"},
		{12*time.Minute + 30*time.Second, "Halfway there! Keep it up."},
		{25 * time.Minute, "Time's up! Recording your hour..."},
	}

	out := &bytes.Buffer{}
	alerter := &SpyScheduleAlerter{}
	p := NewPomodoro(alerter)
	p.Start(out)

	for i, want := range testcases {
		assert.Equal(t, want, alerter.Alerts[i])
	}
	assert.Equal(t, 1, alerter.WaitCalled)
}
