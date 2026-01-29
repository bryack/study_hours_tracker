package cli_test

import (
	"strings"
	"testing"

	"github.com/bryack/study_hours_tracker/adapters/cli"
	"github.com/bryack/study_hours_tracker/testhelpers"
	"github.com/stretchr/testify/assert"
)

func TestCLI(t *testing.T) {

	t.Run("record 'cli' hours", func(t *testing.T) {
		in := strings.NewReader("cli 3")

		store := &testhelpers.StubSubjectStore{
			Hours: map[string]int{},
		}
		trackerCLI := cli.NewCLI(store, in)
		trackerCLI.Run()

		assertRecordTracker(t, store, "cli", 3)
	})

	t.Run("record 'bash' hours", func(t *testing.T) {
		in := strings.NewReader("bash 5")

		store := &testhelpers.StubSubjectStore{
			Hours: map[string]int{},
		}
		trackerCLI := cli.NewCLI(store, in)
		trackerCLI.Run()

		assertRecordTracker(t, store, "bash", 5)
	})
}

func assertRecordTracker(t testing.TB, store *testhelpers.StubSubjectStore, subject string, hours int) {
	t.Helper()

	assert.True(t, len(store.RecordCall) > 0)

	v, ok := store.Hours[subject]
	assert.True(t, ok)
	assert.Equal(t, hours, v)
}
