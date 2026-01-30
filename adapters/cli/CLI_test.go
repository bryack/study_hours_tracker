package cli_test

import (
	"strings"
	"testing"

	"github.com/bryack/study_hours_tracker/adapters/cli"
	"github.com/bryack/study_hours_tracker/testhelpers"
	"github.com/stretchr/testify/assert"
)

func TestCLI(t *testing.T) {
	tests := []struct {
		name            string
		input           string
		expectedSubject string
		expectedHours   int
		expectedErr     error
	}{
		{
			name:            "record 'cli' hours",
			input:           "cli 3",
			expectedSubject: "cli",
			expectedHours:   3,
			expectedErr:     nil,
		},
		{
			name:            "record 'bash' hours",
			input:           "bash 5",
			expectedSubject: "bash",
			expectedHours:   5,
			expectedErr:     nil,
		},
		{
			name:            "handle parsing errors",
			input:           "bufio five",
			expectedSubject: "",
			expectedHours:   0,
			expectedErr:     cli.ErrInvalidHours,
		},
		{
			name:            "not enough arguments",
			input:           "bufio",
			expectedSubject: "",
			expectedHours:   0,
			expectedErr:     cli.ErrNotEnoughArgs,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := strings.NewReader(tt.input)

			store := &testhelpers.StubSubjectStore{
				Hours: map[string]int{},
			}
			trackerCLI := cli.NewCLI(store, in)
			err := trackerCLI.Run()

			assert.ErrorIs(t, err, tt.expectedErr)
			assertRecordTracker(t, store, tt.expectedSubject, tt.expectedHours, tt.expectedErr)
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
