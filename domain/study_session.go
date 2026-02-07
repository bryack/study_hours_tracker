package domain

import "io"

// SessionRunner defines the interface for managing study sessions.
type SessionRunner interface {
	RecordManual(subject string, hours int) error
	RecordPomodoro(subject string, out io.Writer) error
}

// PomodoroRunner represents a timer that can be started for focused study sessions.
type PomodoroRunner interface {
	Start(out io.Writer)
}

// StudySession encapsulates the business logic for recording study hours.
type StudySession struct {
	store          SubjectStore
	pomodoroRunner PomodoroRunner
}

// NewStudySession creates a new study session manager.
func NewStudySession(store SubjectStore, pomodoroRunner PomodoroRunner) *StudySession {
	return &StudySession{
		store:          store,
		pomodoroRunner: pomodoroRunner,
	}
}

// RecordManual records manual study hours.
func (s *StudySession) RecordManual(subject string, hours int) error {
	return s.store.RecordHour(subject, hours)
}

// RecordPomodoro starts a 25-minute Pomodoro session and records it as 1 study hour.
// Note: This is a simplified tracking where 1 Pomodoro = 1 recorded hour for convenience.
func (s *StudySession) RecordPomodoro(subject string, out io.Writer) error {
	s.pomodoroRunner.Start(out)
	return s.store.RecordHour(subject, 1)
}
