package testhelpers

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/bryack/study_hours_tracker/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type StubSubjectStore struct {
	Hours      map[string]int
	RecordCall []string
	Report     domain.Report
	Err        error
}

func (s *StubSubjectStore) RecordHour(subject string, numHours int) error {
	if s.Err != nil {
		return s.Err
	}
	if s.Hours == nil {
		s.Hours = make(map[string]int)
	}
	s.RecordCall = append(s.RecordCall, subject)
	s.Hours[subject] += numHours
	return nil
}

func (s *StubSubjectStore) GetHours(subject string) (int, error) {
	if s.Err != nil {
		return 0, s.Err
	}
	h, ok := s.Hours[subject]
	if !ok {
		return 0, domain.ErrSubjectNotFound
	}
	return h, nil
}

func (s *StubSubjectStore) GetReport() (domain.Report, error) {
	if s.Err != nil {
		return nil, s.Err
	}
	return s.Report, nil
}

func SetupTestContainer(t testing.TB) string {
	t.Helper()
	ctx := context.Background()

	user := "testuser"
	pass := "testpass"
	dbName := "testDB"

	req := testcontainers.ContainerRequest{
		Image:        "postgres:15-alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     user,
			"POSTGRES_PASSWORD": pass,
			"POSTGRES_DB":       dbName,
		},
		WaitingFor: wait.ForLog("database system is ready to accept connections").WithOccurrence(2).WithStartupTimeout(10 * time.Second),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(t, err)

	t.Cleanup(func() {
		assert.NoError(t, container.Terminate(ctx))
	})

	host, err := container.Host(ctx)
	require.NoError(t, err, "failed to get host for container")
	port, err := container.MappedPort(ctx, "5432")
	require.NoError(t, err, "failed to get port for container")

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, pass, host, port.Port(), dbName)
}
