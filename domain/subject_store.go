package domain

type SubjectStore interface {
	GetHours(subject string) (int, error)
	RecordHour(subject string, numHours int) error
	GetReport() (Report, error)
}
