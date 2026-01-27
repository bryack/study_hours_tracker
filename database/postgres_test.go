package database

import (
	"testing"

	"github.com/bryack/study_hours_tracker/testhelpers"
	"github.com/stretchr/testify/assert"
)

func TestRecordAndGetHours(t *testing.T) {
	connStr := testhelpers.SetupTestContainer(t)
	store, err := NewPostgresSubjectStore(connStr)
	if err != nil {
		t.Fatal(err)
	}
	t.Run("should_successfully_connect_and_create_table", func(t *testing.T) {
		store := &PostgresSubjectStore{}
		err := store.initDatabase(connStr)

		t.Cleanup(func() {
			if store.db != nil {
				store.db.Close()
			}
		})

		assert.NoError(t, err)
		err = store.db.Ping()
		assert.NoError(t, err, "Database must be reachable via Ping")

		err = store.createTable()
		var tableName string
		query := `SELECT table_name FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'subjects'`

		err = store.db.QueryRow(query).Scan(&tableName)
		assert.NoError(t, err, "Table 'subjects' should exist")
		assert.Equal(t, "subjects", tableName)
	})

	t.Run("record hours for tdd", func(t *testing.T) {
		err := store.RecordHour("tdd", 2)
		assert.NoError(t, err)

		err = store.RecordHour("tdd", 3)
		assert.NoError(t, err)

		var count int
		err = store.db.QueryRow("SELECT COUNT(*) FROM subjects").Scan(&count)
		assert.NoError(t, err)
		assert.True(t, count == 1, "Rows count should still be 1")

		var hours int
		err = store.db.QueryRow("SELECT hours FROM subjects WHERE subject = 'tdd'").Scan(&hours)
		assert.NoError(t, err)
		assert.Equal(t, 5, hours, "Hours should be summed up (2 + 3)")
	})

	t.Run("get hours for tdd", func(t *testing.T) {
		h, err := store.GetHours("tdd")
		assert.NoError(t, err)

		var hours int
		err = store.db.QueryRow("SELECT hours FROM subjects WHERE subject = 'tdd'").Scan(&hours)
		assert.NoError(t, err)
		assert.Equal(t, h, hours)
	})

	t.Run("get hours for nonexistent subject", func(t *testing.T) {
		h, err := store.GetHours("nonexistent")
		assert.ErrorIs(t, err, ErrSubjectNotFound)
		assert.Equal(t, 0, h)
	})
}
