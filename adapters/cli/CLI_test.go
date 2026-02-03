package cli_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/bryack/study_hours_tracker/adapters/cli"
	"github.com/bryack/study_hours_tracker/testhelpers"
	"github.com/stretchr/testify/assert"
)

type SpyPomodoroRunner struct {
	StartCallCount int
}

func (s *SpyPomodoroRunner) Start() {
	s.StartCallCount++
}

func TestCLI(t *testing.T) {
	tests := []struct {
		name              string
		input             string
		expectedOut       string
		expectedSubject   string
		expectedHours     int
		shouldRecord      bool
		expectedSpyCalled int
	}{
		{
			name:              "record 'cli' hours",
			input:             "cli 3",
			expectedOut:       cli.GreetingString,
			expectedSubject:   "cli",
			expectedHours:     3,
			shouldRecord:      true,
			expectedSpyCalled: 0,
		},
		{
			name:              "handle parsing errors",
			input:             "bufio five",
			expectedOut:       cli.GreetingString + "\nfailed to extract subject and hours",
			expectedSubject:   "",
			expectedHours:     0,
			shouldRecord:      false,
			expectedSpyCalled: 0,
		},
		{
			name:              "not enough arguments",
			input:             "bufio",
			expectedOut:       cli.GreetingString + "\nfailed to extract subject and hours",
			expectedSubject:   "",
			expectedHours:     0,
			shouldRecord:      false,
			expectedSpyCalled: 0,
		},
		{
			name:              "negative number of hours",
			input:             "bufio -2",
			expectedOut:       cli.GreetingString + "\nfailed to extract subject and hours",
			expectedSubject:   "",
			expectedHours:     0,
			shouldRecord:      false,
			expectedSpyCalled: 0,
		},
		{
			name:              "record multiple sessions",
			input:             "cli 3\nbash 2",
			expectedOut:       cli.GreetingString,
			expectedSubject:   "bash",
			expectedHours:     2,
			shouldRecord:      true,
			expectedSpyCalled: 0,
		},
		{
			name:              "continue after error",
			input:             "cli 3\ninvalid_data\nbash 2",
			expectedOut:       cli.GreetingString + "\nfailed to extract subject and hours",
			expectedSubject:   "bash",
			expectedHours:     2,
			shouldRecord:      true,
			expectedSpyCalled: 0,
		},
		{
			name:              "start pomodoro for tdd",
			input:             "pomodoro tdd",
			expectedOut:       cli.GreetingString + "\nPomodoro started...",
			expectedSubject:   "tdd",
			expectedHours:     1,
			shouldRecord:      true,
			expectedSpyCalled: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			in := strings.NewReader(tt.input)
			out := &bytes.Buffer{}

			store := &testhelpers.StubSubjectStore{
				Hours:      map[string]int{},
				RecordCall: []string{},
			}

			pomodoroSpy := &SpyPomodoroRunner{}
			trackerCLI := cli.NewCLI(store, in, out, pomodoroSpy)
			err := trackerCLI.Run()

			assert.NoError(t, err)
			assertRecordTracker(t, store, tt.expectedSubject, tt.expectedHours, tt.shouldRecord)
			assert.Contains(t, strings.TrimSpace(out.String()), tt.expectedOut)
			assert.Equal(t, tt.expectedSpyCalled, pomodoroSpy.StartCallCount)
		})
	}
}

func assertRecordTracker(t testing.TB, store *testhelpers.StubSubjectStore, subject string, hours int, shouldRecord bool) {
	t.Helper()

	if shouldRecord {
		assert.True(t, len(store.RecordCall) > 0, "expected RecordHour to be called, but it wasn't")

		got, ok := store.Hours[subject]
		assert.True(t, ok, "expected subject %q to be recorded, but it wasn't found", subject)
		assert.Equal(t, got, hours, "expected %d hours for %q, got %d", hours, subject, got)
	} else {
		assert.Equal(t, 0, len(store.RecordCall), "expected no calls to RecordHour, but got %d calls", len(store.RecordCall))
	}
}
