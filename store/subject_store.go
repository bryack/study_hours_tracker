package store

import "github.com/bryack/study_hours_tracker/domain"

type SubjectStore interface {
	GetHours(subject string) (int, error)
	RecordHour(subject string, numHours int) error
	GetReport() (domain.Report, error)
}
