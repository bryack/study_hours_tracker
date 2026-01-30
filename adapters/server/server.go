package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/bryack/study_hours_tracker/domain"
	"github.com/bryack/study_hours_tracker/store"
)

const (
	jsonContentType = "application/json"
	reportPath      = "/report"
	trackerPath     = "/tracker/"
)

type StudyServer struct {
	Store store.SubjectStore
	http.Handler
}

func NewStudyServer(store store.SubjectStore) *StudyServer {
	s := &StudyServer{}

	s.Store = store

	router := http.NewServeMux()
	router.Handle(reportPath, http.HandlerFunc(s.reportHandler))
	router.Handle(trackerPath, http.HandlerFunc(s.trackerHandler))

	s.Handler = router

	return s
}

func (s *StudyServer) reportHandler(w http.ResponseWriter, r *http.Request) {
	studyActivities, err := s.Store.GetReport()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("content-type", jsonContentType)
	if err := json.NewEncoder(w).Encode(studyActivities); err != nil {
		log.Println("failed to encode:", err)
	}
}

func (s *StudyServer) trackerHandler(w http.ResponseWriter, r *http.Request) {
	subject := strings.TrimPrefix(r.URL.Path, trackerPath)

	if strings.TrimSpace(subject) == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

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
		if errors.Is(err, domain.ErrSubjectNotFound) {
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
