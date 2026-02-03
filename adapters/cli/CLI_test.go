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
	SpyCalled int
}

func (s *SpyPomodoroRunner) Start() {
	s.SpyCalled++
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
			expectedOut:       cli.GretingString,
			expectedSubject:   "cli",
			expectedHours:     3,
			shouldRecord:      true,
			expectedSpyCalled: 0,
		},
		{
			name:              "handle parsing errors",
			input:             "bufio five",
			expectedOut:       cli.GretingString + "\nfailed to extract subject and hours",
			expectedSubject:   "",
			expectedHours:     0,
			shouldRecord:      false,
			expectedSpyCalled: 0,
		},
		{
			name:              "not enough arguments",
			input:             "bufio",
			expectedOut:       cli.GretingString + "\nfailed to extract subject and hours",
			expectedSubject:   "",
			expectedHours:     0,
			shouldRecord:      false,
			expectedSpyCalled: 0,
		},
		{
			name:              "negative number of hours",
			input:             "bufio -2",
			expectedOut:       cli.GretingString + "\nfailed to extract subject and hours",
			expectedSubject:   "",
			expectedHours:     0,
			shouldRecord:      false,
			expectedSpyCalled: 0,
		},
		{
			name:              "record multiple sessions",
			input:             "cli 3\nbash 2",
			expectedOut:       cli.GretingString,
			expectedSubject:   "bash",
			expectedHours:     2,
			shouldRecord:      true,
			expectedSpyCalled: 0,
		},
		{
			name:              "continue after error",
			input:             "cli 3\ninvalid_data\nbash 2",
			expectedOut:       cli.GretingString + "\nfailed to extract subject and hours",
			expectedSubject:   "bash",
			expectedHours:     2,
			shouldRecord:      true,
			expectedSpyCalled: 0,
		},
		{
			name:              "start pomodoro for tdd",
			input:             "pomodoro tdd",
			expectedOut:       cli.GretingString + "\nPomodoro started...",
			expectedSubject:   "tdd",
			expectedHours:     1,
			shouldRecord:      true,
			expectedSpyCalled: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
			assert.Equal(t, tt.expectedSpyCalled, pomodoroSpy.SpyCalled)
		})
	}
}

func assertRecordTracker(t testing.TB, store *testhelpers.StubSubjectStore, subject string, hours int, shouldRecord bool) {
	t.Helper()

	if shouldRecord {
		assert.True(t, len(store.RecordCall) > 0)

		v, ok := store.Hours[subject]
		assert.True(t, ok)
		assert.Equal(t, hours, v)
	} else {
		assert.Equal(t, 0, len(store.RecordCall))
	}
}
