package server

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/bryack/study_hours_tracker/database"
)

type SubjectStore interface {
	GetHours(subject string) (int, error)
	RecordHour(subject string, numHours int) error
}

type StudyServer struct {
	Store SubjectStore
}

func (s *StudyServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	subject := strings.TrimPrefix(r.URL.Path, "/tracker/")

	switch r.Method {
	case http.MethodPost:
		s.processPostRequest(w, r, subject)
	case http.MethodGet:
		s.processGetRequest(w, subject)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (s *StudyServer) processGetRequest(w http.ResponseWriter, subject string) {
	hours, err := s.Store.GetHours(subject)
	if err != nil {
		if errors.Is(err, database.ErrSubjectNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, hours)
}

func (s *StudyServer) processPostRequest(w http.ResponseWriter, r *http.Request, subject string) {
	h, err := strconv.Atoi(r.URL.Query().Get("hours"))
	if err != nil || h <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = s.Store.RecordHour(subject, h)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}
