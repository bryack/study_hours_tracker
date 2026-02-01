package cli_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/bryack/study_hours_tracker/adapters/cli"
	"github.com/bryack/study_hours_tracker/testhelpers"
	"github.com/stretchr/testify/assert"
)

type SpySleeper struct {
	DurationSlept time.Duration
}

func (s *SpySleeper) Sleep(duration time.Duration) {
	s.DurationSlept = duration
}

func TestCLI(t *testing.T) {
	tests := []struct {
		name            string
		input           string
		expectedOut     string
		expectedSubject string
		expectedHours   int
		expectedSleep   time.Duration
		expectedErr     error
	}{
		{
			name:            "record 'cli' hours",
			input:           "cli 3",
			expectedOut:     "",
			expectedSubject: "cli",
			expectedHours:   3,
			expectedSleep:   0,
			expectedErr:     nil,
		},
		{
			name:            "record 'bash' hours",
			input:           "bash 5",
			expectedOut:     "",
			expectedSubject: "bash",
			expectedHours:   5,
			expectedSleep:   0,
			expectedErr:     nil,
		},
		{
			name:            "handle parsing errors",
			input:           "bufio five",
			expectedOut:     "",
			expectedSubject: "",
			expectedHours:   0,
			expectedSleep:   0,
			expectedErr:     cli.ErrInvalidHours,
		},
		{
			name:            "not enough arguments",
			input:           "bufio",
			expectedOut:     "",
			expectedSubject: "",
			expectedHours:   0,
			expectedSleep:   0,
			expectedErr:     cli.ErrNotEnoughArgs,
		},
		{
			name:            "negative number of hours",
			input:           "bufio -2",
			expectedOut:     "",
			expectedSubject: "",
			expectedHours:   0,
			expectedSleep:   0,
			expectedErr:     cli.ErrInvalidHours,
		},
		{
			name:            "start pomodoro for tdd",
			input:           "pomodoro tdd",
			expectedOut:     "Pomodoro started...",
			expectedSubject: "tdd",
			expectedHours:   1,
			expectedSleep:   25 * time.Minute,
			expectedErr:     nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := strings.NewReader(tt.input)
			out := &bytes.Buffer{}

			store := &testhelpers.StubSubjectStore{
				Hours: map[string]int{},
			}
			sleeper := &SpySleeper{}
			trackerCLI := cli.NewCLI(store, in, out, sleeper)
			err := trackerCLI.Run()

			assert.ErrorIs(t, err, tt.expectedErr)
			assertRecordTracker(t, store, tt.expectedSubject, tt.expectedHours, tt.expectedErr)
			assert.Equal(t, tt.expectedOut, strings.TrimSpace(out.String()))
			assert.Equal(t, tt.expectedSleep, sleeper.DurationSlept)
		})
	}
}

func assertRecordTracker(t testing.TB, store *testhelpers.StubSubjectStore, subject string, hours int, err error) {
	t.Helper()

	if err == nil {
		assert.True(t, len(store.RecordCall) > 0)

		v, ok := store.Hours[subject]
		assert.True(t, ok)
		assert.Equal(t, hours, v)
	} else {
		assert.Equal(t, 0, len(store.RecordCall))
	}
}
