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

func (s *SpySleeper) Wait(duration time.Duration) {
	s.DurationSlept = duration
}

func TestCLI(t *testing.T) {
	tests := []struct {
		name            string
		input           string
		expectedOut     string
		expectedSubject string
		expectedHours   int
		expectedCalls   []string
		expectedMap     map[string]int
		expectedSleep   time.Duration
		shouldRecord    bool
		expectedErr     error
	}{
		{
			name:            "record 'cli' hours",
			input:           "cli 3",
			expectedOut:     cli.GretingString,
			expectedSubject: "cli",
			expectedHours:   3,
			expectedCalls:   []string{"cli"},
			expectedSleep:   0,
			shouldRecord:    true,
			expectedErr:     nil,
		},
		{
			name:            "record 'bash' hours",
			input:           "bash 5",
			expectedOut:     cli.GretingString,
			expectedSubject: "bash",
			expectedHours:   5,
			expectedCalls:   []string{"bash"},
			expectedSleep:   0,
			shouldRecord:    true,
			expectedErr:     nil,
		},
		{
			name:            "handle parsing errors",
			input:           "bufio five",
			expectedOut:     cli.GretingString + "\nfailed to extract subject and hours: failed to parse hours five: strconv.Atoi: parsing \"five\": invalid syntax",
			expectedSubject: "",
			expectedHours:   0,
			expectedCalls:   []string{},
			expectedSleep:   0,
			shouldRecord:    false,
			expectedErr:     nil,
		},
		{
			name:            "not enough arguments",
			input:           "bufio",
			expectedOut:     cli.GretingString + "\nfailed to extract subject and hours: failed to parse: should be 2 arguments, got: 1",
			expectedSubject: "",
			expectedHours:   0,
			expectedCalls:   []string{},
			expectedSleep:   0,
			shouldRecord:    false,
			expectedErr:     nil,
		},
		{
			name:            "negative number of hours",
			input:           "bufio -2",
			expectedOut:     cli.GretingString + "\nfailed to extract subject and hours: failed to parse hours -2, should be 1 or more: <nil>",
			expectedSubject: "",
			expectedHours:   0,
			expectedCalls:   []string{},
			expectedSleep:   0,
			shouldRecord:    false,
			expectedErr:     nil,
		},
		{
			name:            "start pomodoro for tdd",
			input:           "pomodoro tdd",
			expectedOut:     cli.GretingString + "\nPomodoro started...",
			expectedSubject: "tdd",
			expectedHours:   1,
			expectedCalls:   []string{"tdd"},
			expectedSleep:   25 * time.Minute,
			shouldRecord:    true,
			expectedErr:     nil,
		},
		{
			name:            "record multiple sessions",
			input:           "cli 3\nbash 2",
			expectedOut:     cli.GretingString,
			expectedSubject: "bash",
			expectedHours:   2,
			expectedCalls:   []string{"cli", "bash"},
			expectedMap:     map[string]int{"cli": 3, "bash": 2},
			expectedSleep:   0,
			shouldRecord:    true,
			expectedErr:     nil,
		},
		{
			name:            "continue after error",
			input:           "cli 3\ninvalid_data\nbash 2",
			expectedOut:     cli.GretingString + "\nfailed to extract subject and hours: failed to parse: should be 2 arguments, got: 1",
			expectedSubject: "bash",
			expectedHours:   2,
			expectedCalls:   []string{"cli", "bash"},
			expectedMap:     map[string]int{"cli": 3, "bash": 2},
			expectedSleep:   0,
			shouldRecord:    true,
			expectedErr:     nil,
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
			sleeper := &SpySleeper{}
			trackerCLI := cli.NewCLI(store, in, out, sleeper)
			err := trackerCLI.Run()

			assert.ErrorIs(t, err, tt.expectedErr)
			assertRecordTracker(t, store, tt.expectedSubject, tt.expectedHours, tt.shouldRecord)
			assert.Equal(t, tt.expectedCalls, store.RecordCall)
			if tt.expectedMap != nil {
				assert.Equal(t, tt.expectedMap, store.Hours)
			}
			assert.Equal(t, tt.expectedOut, strings.TrimSpace(out.String()))
			assert.Equal(t, tt.expectedSleep, sleeper.DurationSlept)
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
