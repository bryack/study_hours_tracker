package server

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/bryack/study_hours_tracker/adapters/database"
	"github.com/bryack/study_hours_tracker/domain"
	"github.com/bryack/study_hours_tracker/testhelpers"
	"github.com/gorilla/websocket"

	"github.com/stretchr/testify/assert"
)

func TestGETSubjects(t *testing.T) {
	store := &testhelpers.StubSubjectStore{
		Hours: map[string]int{
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
		failedStore := &testhelpers.StubSubjectStore{
			GetHoursErr: errors.New("database connection lost"),
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
		recordHourErr error
	}{
		{
			name:          "record TDD hours as positive number",
			path:          "/tracker/tdd?hours=5",
			subjectsSlice: []string{"tdd"},
			numHours:      5,
			expectedCode:  202,
			recordHourErr: nil,
		},
		{
			name:          "record http hours as string",
			path:          "/tracker/tdd?hours=aaa",
			subjectsSlice: []string{},
			numHours:      0,
			expectedCode:  400,
			recordHourErr: nil,
		},
		{
			name:          "record http hours as negative number",
			path:          "/tracker/tdd?hours=-1",
			subjectsSlice: []string{},
			numHours:      0,
			expectedCode:  400,
			recordHourErr: nil,
		},
		{
			name:          "expected 500",
			path:          "/tracker/db?hours=2",
			subjectsSlice: []string{},
			numHours:      0,
			expectedCode:  500,
			recordHourErr: errors.New("persistent storage failure"),
		},
		{
			name:          "empty subject",
			path:          "/tracker/?hours=2",
			subjectsSlice: []string{},
			numHours:      0,
			expectedCode:  400,
			recordHourErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := &testhelpers.StubSubjectStore{
				Hours:         map[string]int{},
				RecordCall:    []string{},
				RecordHourErr: tt.recordHourErr,
			}
			server := NewStudyServer(store)
			request, err := http.NewRequest(http.MethodPost, tt.path, nil)
			if err != nil {
				t.Fatal(err)
			}
			response := httptest.NewRecorder()
			server.ServeHTTP(response, request)

			assert.Equal(t, tt.expectedCode, response.Code)
			assert.Equal(t, tt.subjectsSlice, store.RecordCall)
			assert.Equal(t, tt.numHours, store.Hours["tdd"])
		})
	}
}

func TestMethodNotAllowed(t *testing.T) {
	store := &testhelpers.StubSubjectStore{
		Hours: map[string]int{},
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
	svr := NewStudyServer(store)

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
	h, err := store.GetHours("tdd")
	assert.NoError(t, err)
	assert.Equal(t, strconv.Itoa(h), response.Body.String())
}

func TestReport(t *testing.T) {
	t.Run("returns 200 on /report", func(t *testing.T) {
		wantedReport := domain.Report{
			{Subject: "Docker", Hours: 4},
			{Subject: "TDD", Hours: 6},
		}
		store := &testhelpers.StubSubjectStore{
			Report: wantedReport,
		}
		server := NewStudyServer(store)
		request, err := http.NewRequest(http.MethodGet, "/report", nil)
		assert.NoError(t, err)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := getReportFromResponse(t, response.Body)

		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, wantedReport, got)
		assert.Equal(t, jsonContentType, response.Result().Header.Get("content-type"))

	})
	t.Run("handle 500", func(t *testing.T) {
		store := &testhelpers.StubSubjectStore{
			Hours:        map[string]int{},
			GetReportErr: errors.New("database connection failed"),
		}
		server := NewStudyServer(store)

		request, err := http.NewRequest(http.MethodGet, "/report", nil)
		assert.NoError(t, err)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assert.Equal(t, http.StatusInternalServerError, response.Code)
	})
}

func getReportFromResponse(t testing.TB, body io.Reader) domain.Report {
	t.Helper()
	var report domain.Report
	err := json.NewDecoder(body).Decode(&report)
	if err != nil {
		t.Fatalf("Unable to parse response from server %q into slice of StudyActivity, '%v'", body, err)
	}
	return report
}

func TestStudy(t *testing.T) {

	t.Run("GET /study returns 200", func(t *testing.T) {
		store := &testhelpers.StubSubjectStore{}
		server := NewStudyServer(store)
		request := newStudyRequest(t)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assert.Equal(t, http.StatusOK, response.Code)
	})
	t.Run("upgrade to websocket", func(t *testing.T) {
		store := &testhelpers.StubSubjectStore{}
		subject := "websockets"
		server := httptest.NewServer(NewStudyServer(store))
		defer server.Close()

		wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"
		conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
		if err != nil {
			t.Fatalf("failed to open a ws connection on %q: %v", wsURL, err)
		}
		defer conn.Close()

		if err := conn.WriteMessage(websocket.TextMessage, []byte(subject)); err != nil {
			t.Fatalf("failed to send message %q ovew ws connection: %v", subject, err)
		}

		time.Sleep(10 * time.Millisecond)
		assert.Equal(t, store.RecordCall[0], subject)
	})
}

func newStudyRequest(t *testing.T) *http.Request {
	request, err := http.NewRequest(http.MethodGet, "/study", nil)
	assert.NoError(t, err)
	return request
}
