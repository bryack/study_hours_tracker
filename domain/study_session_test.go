package domain_test

import (
	"errors"
	"testing"

	"github.com/bryack/study_hours_tracker/domain"
	"github.com/bryack/study_hours_tracker/testhelpers"
	"github.com/stretchr/testify/assert"
)

type SpyPomodoroRunner struct {
	StartCallCount int
}

func (s *SpyPomodoroRunner) Start() {
	s.StartCallCount++
}

func TestStudySession_RecordPomodoro(t *testing.T) {
	t.Run("starts pomodoro and records 1 hour", func(t *testing.T) {
		store := &testhelpers.StubSubjectStore{
			Hours:      map[string]int{},
			RecordCall: []string{},
		}

		pomodoroSpy := &SpyPomodoroRunner{}
		session := domain.NewStudySession(store, pomodoroSpy)

		err := session.RecordPomodoro("cli")
		assert.NoError(t, err)

		v, ok := store.Hours["cli"]
		assert.True(t, ok)
		assert.Equal(t, 1, v, "should record 1 hour")
		assert.Equal(t, 1, pomodoroSpy.StartCallCount, "should start pomodoro once")
	})
	t.Run("returns error if store fails", func(t *testing.T) {
		store := &testhelpers.StubSubjectStore{
			Hours:         map[string]int{},
			RecordCall:    []string{},
			RecordHourErr: errors.New("persistent storage failure"),
		}

		pomodoroSpy := &SpyPomodoroRunner{}
		session := domain.NewStudySession(store, pomodoroSpy)

		err := session.RecordPomodoro("cli")
		assert.Error(t, err)

		v, ok := store.Hours["cli"]
		assert.True(t, !ok)
		assert.Equal(t, 0, v, "should not record 1 hour")
		assert.Equal(t, 1, pomodoroSpy.StartCallCount, "should still start pomodoro")
	})
}

func TestStudySession_RecordManual(t *testing.T) {
	t.Run("records manual hours to store", func(t *testing.T) {
		store := &testhelpers.StubSubjectStore{
			Hours:      map[string]int{},
			RecordCall: []string{},
		}

		pomodoroSpy := &SpyPomodoroRunner{}
		session := domain.NewStudySession(store, pomodoroSpy)

		err := session.RecordManual("cli", 3)
		assert.NoError(t, err)

		v, ok := store.Hours["cli"]
		assert.True(t, ok)
		assert.Equal(t, 3, v)
		assert.Equal(t, 0, pomodoroSpy.StartCallCount, "should not start pomodoro")
	})
}
