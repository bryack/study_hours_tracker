package cli_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/bryack/study_hours_tracker/adapters/cli"
	"github.com/bryack/study_hours_tracker/testhelpers"
	"github.com/stretchr/testify/assert"
)

func TestCLI(t *testing.T) {
	tests := []struct {
		name                  string
		input                 string
		expectedOut           string
		expectedManualCalls   map[string]int
		expectedPomodoroCalls []string
	}{
		{
			name:                  "record 'cli' hours",
			input:                 "cli 3",
			expectedOut:           cli.GreetingString,
			expectedManualCalls:   map[string]int{"cli": 3},
			expectedPomodoroCalls: []string{},
		},
		{
			name:                  "handle parsing errors",
			input:                 "bufio five",
			expectedOut:           cli.GreetingString + "\nfailed to extract subject and hours",
			expectedManualCalls:   map[string]int{},
			expectedPomodoroCalls: []string{},
		},
		{
			name:                  "not enough arguments",
			input:                 "bufio",
			expectedOut:           cli.GreetingString + "\nfailed to extract subject and hours",
			expectedManualCalls:   map[string]int{},
			expectedPomodoroCalls: []string{},
		},
		{
			name:                  "negative number of hours",
			input:                 "bufio -2",
			expectedOut:           cli.GreetingString + "\nfailed to extract subject and hours",
			expectedManualCalls:   map[string]int{},
			expectedPomodoroCalls: []string{},
		},
		{
			name:                  "record multiple sessions",
			input:                 "cli 3\nbash 2",
			expectedOut:           cli.GreetingString,
			expectedManualCalls:   map[string]int{"cli": 3, "bash": 2},
			expectedPomodoroCalls: []string{},
		},
		{
			name:                  "continue after error",
			input:                 "cli 3\ninvalid_data\nbash 2",
			expectedOut:           cli.GreetingString + "\nfailed to extract subject and hours",
			expectedManualCalls:   map[string]int{"cli": 3, "bash": 2},
			expectedPomodoroCalls: []string{},
		},
		{
			name:                  "start pomodoro for tdd",
			input:                 "pomodoro tdd",
			expectedOut:           cli.GreetingString + "\nPomodoro started...",
			expectedManualCalls:   map[string]int{},
			expectedPomodoroCalls: []string{"tdd"},
		},
		{
			name:                  "quit command exits gracefully",
			input:                 "quit",
			expectedOut:           cli.GreetingString + "\nGoodbye!",
			expectedManualCalls:   map[string]int{},
			expectedPomodoroCalls: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			session := &testhelpers.SpySession{
				ManualCalls:   map[string]int{},
				PomodoroCalls: []string{},
			}
			in := strings.NewReader(tt.input)
			out := &bytes.Buffer{}

			trackerCLI := cli.NewCLI(in, out, session)
			err := trackerCLI.Run()
			assert.NoError(t, err)

			assert.Equal(t, tt.expectedManualCalls, session.ManualCalls)
			assert.Equal(t, tt.expectedPomodoroCalls, session.PomodoroCalls)
			assert.Contains(t, strings.TrimSpace(out.String()), tt.expectedOut)
		})
	}
}
