package server

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/bryack/study_hours_tracker/database"
	"github.com/bryack/study_hours_tracker/testhelpers"

	"github.com/stretchr/testify/assert"
)

type StubSubjectStore struct {
	hours      map[string]int
	recordCall []string
	err        error
}

func (s *StubSubjectStore) RecordHour(subject string, numHours int) error {
	if s.err != nil {
		return s.err
	}
	s.recordCall = append(s.recordCall, subject)
	s.hours[subject] += numHours
	return nil
}

func (s *StubSubjectStore) GetHours(subject string) (int, error) {
	if s.err != nil {
		return 0, s.err
	}
	h, ok := s.hours[subject]
	if !ok {
		return 0, database.ErrSubjectNotFound
	}
	return h, nil
}

func TestGETSubjects(t *testing.T) {
	store := &StubSubjectStore{
		hours: map[string]int{
			"tdd":  20,
			"http": 10,
		},
	}
	server := NewStudyServer(store)
	t.Run("returns TDD hours", func(t *testing.T) {
		tddHours := "20"
		request, err := http.NewRequest(http.MethodGet, "/tracker/tdd", nil)
		if err != nil {
			t.Fatal(err)
		}
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assert.Equal(t, response.Code, http.StatusOK)
		assert.Equal(t, response.Body.String(), tddHours)
	})
	t.Run("returns http hours", func(t *testing.T) {
		httpHours := "10"
		request, err := http.NewRequest(http.MethodGet, "/tracker/http", nil)
		if err != nil {
			t.Fatal(err)
		}
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assert.Equal(t, response.Code, http.StatusOK)
		assert.Equal(t, response.Body.String(), httpHours)
	})

	t.Run("handle 404", func(t *testing.T) {
		request, err := http.NewRequest(http.MethodGet, "/tracker/java", nil)
		if err != nil {
			t.Fatal(err)
		}
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assert.Equal(t, response.Code, http.StatusNotFound)
	})

	t.Run("returns 500 when store fails", func(t *testing.T) {
		failedStore := &StubSubjectStore{
			err: errors.New("database connection lost"),
		}
		failedServer := NewStudyServer(failedStore)
		request, err := http.NewRequest(http.MethodGet, "/tracker/tdd", nil)
		assert.NoError(t, err)
		response := httptest.NewRecorder()

		failedServer.ServeHTTP(response, request)

		assert.Equal(t, http.StatusInternalServerError, response.Code)
	})
}

func TestPostHoursToSubject(t *testing.T) {
	tests := []struct {
		name          string
		path          string
		subjectsSlice []string
		numHours      int
		expectedCode  int
		expectedErr   error
	}{
		{
			name:          "record TDD hours as positive number",
			path:          "/tracker/tdd?hours=5",
			subjectsSlice: []string{"tdd"},
			numHours:      5,
			expectedCode:  202,
			expectedErr:   nil,
		},
		{
			name:          "record http hours as string",
			path:          "/tracker/tdd?hours=aaa",
			subjectsSlice: []string{},
			numHours:      0,
			expectedCode:  400,
			expectedErr:   nil,
		},
		{
			name:          "record http hours as negative number",
			path:          "/tracker/tdd?hours=-1",
			subjectsSlice: []string{},
			numHours:      0,
			expectedCode:  400,
			expectedErr:   nil,
		},
		{
			name:          "expected 500",
			path:          "/tracker/db?hours=2",
			subjectsSlice: []string{},
			numHours:      0,
			expectedCode:  500,
			expectedErr:   errors.New("persistent storage failure"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := &StubSubjectStore{
				hours:      map[string]int{},
				recordCall: []string{},
				err:        tt.expectedErr,
			}
			server := NewStudyServer(store)
			request, err := http.NewRequest(http.MethodPost, tt.path, nil)
			if err != nil {
				t.Fatal(err)
			}
			response := httptest.NewRecorder()
			server.ServeHTTP(response, request)

			assert.Equal(t, tt.expectedCode, response.Code)
			assert.Equal(t, tt.subjectsSlice, store.recordCall)
			assert.Equal(t, tt.numHours, store.hours["tdd"])
		})
	}
}

func TestMethodNotAllowed(t *testing.T) {
	store := &StubSubjectStore{
		hours: map[string]int{},
	}
	server := NewStudyServer(store)
	t.Run("handle 405", func(t *testing.T) {
		request, err := http.NewRequest(http.MethodPut, "/tracker/tdd", nil)
		if err != nil {
			t.Fatal(err)
		}
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assert.Equal(t, response.Code, http.StatusMethodNotAllowed)
	})
}

func TestRacePostgresSubjectStore(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping race condition test in short mode")
	}
	connStr := testhelpers.SetupTestContainer(t)
	store, err := database.NewPostgresSubjectStore(connStr)
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}
	svr := NewStudyServer(store)

	const concurrentRequests = 100
	const hoursPerRequest = 2

	var wg sync.WaitGroup

	for range concurrentRequests {
		wg.Add(1)
		go func() {
			defer wg.Done()
			req, err := http.NewRequest(http.MethodPost, "/tracker/tdd?hours=2", nil)
			assert.NoError(t, err)
			svr.ServeHTTP(httptest.NewRecorder(), req)
		}()
	}
	wg.Wait()

	got, err := store.GetHours("tdd")
	assert.NoError(t, err)

	expected := concurrentRequests * hoursPerRequest
	assert.Equal(t, expected, got, "Race condition detected: sums do not match")
}

func TestRecordingHoursAndRetrievingThem(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}
	connStr := testhelpers.SetupTestContainer(t)
	store, err := database.NewPostgresSubjectStore(connStr)
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}
	svr := &StudyServer{
		Store: store,
	}

	postReq, err := http.NewRequest(http.MethodPost, "/tracker/tdd?hours=1", nil)
	if err != nil {
		t.Fatal(err)
	}
	getReq, err := http.NewRequest(http.MethodGet, "/tracker/tdd", nil)
	if err != nil {
		t.Fatal(err)
	}
	response := httptest.NewRecorder()

	svr.ServeHTTP(httptest.NewRecorder(), postReq)
	svr.ServeHTTP(httptest.NewRecorder(), postReq)
	svr.ServeHTTP(response, getReq)

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, "2", response.Body.String())
}

func TestReport(t *testing.T) {
	store := &StubSubjectStore{
		hours: map[string]int{},
	}
	server := NewStudyServer(store)
	t.Run("returns 200 on /report", func(t *testing.T) {
		request, err := http.NewRequest(http.MethodGet, "/report", nil)
		assert.NoError(t, err)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assert.Equal(t, http.StatusOK, response.Code)
	})
}
