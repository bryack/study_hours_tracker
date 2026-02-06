package domain

// PomodoroRunner represents a timer that can be started for focused study sessions.
type PomodoroRunner interface {
	Start()
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

// RecordPomodoro starts a Pomodoro session and records 1 hour.
func (s *StudySession) RecordPomodoro(subject string) error {
	s.pomodoroRunner.Start()
	return s.store.RecordHour(subject, 1)
}
