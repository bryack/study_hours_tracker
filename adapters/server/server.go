package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/bryack/study_hours_tracker/domain"
	"github.com/gorilla/websocket"
)

const (
	jsonContentType  = "application/json"
	reportPath       = "/report"
	trackerPath      = "/tracker/"
	studyPath        = "/study"
	websocketPath    = "/ws"
	htmlTemplatePath = "../../study.html"
)

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type wsMessage struct {
	Command string `json:"command"`
	Subject string `json:"subject"`
	Hours   int    `json:"hours,omitempty"` // Optional, only for record_manual
}

type StudyServer struct {
	store    domain.SubjectStore
	template *template.Template
	session  domain.SessionRunner
	http.Handler
}

func NewStudyServer(store domain.SubjectStore, session domain.SessionRunner) (*StudyServer, error) {
	s := &StudyServer{}

	tmpl, err := template.ParseFiles(htmlTemplatePath)
	if err != nil {
		return nil, fmt.Errorf("failed to load template %q: %s", htmlTemplatePath, err)
	}

	s.store = store
	s.template = tmpl
	s.session = session

	router := http.NewServeMux()
	router.Handle(reportPath, http.HandlerFunc(s.reportHandler))
	router.Handle(trackerPath, http.HandlerFunc(s.trackerHandler))
	router.Handle(studyPath, http.HandlerFunc(s.studyHandler))
	router.Handle(websocketPath, http.HandlerFunc(s.webSocketHandler))

	s.Handler = router

	return s, nil
}

func (s *StudyServer) reportHandler(w http.ResponseWriter, r *http.Request) {
	studyActivities, err := s.store.GetReport()
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

func (s *StudyServer) studyHandler(w http.ResponseWriter, r *http.Request) {
	s.template.Execute(w, nil)
}

func (s *StudyServer) webSocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, _ := wsUpgrader.Upgrade(w, r, nil)

	for {
		_, msgBytes, _ := conn.ReadMessage()
		var msg wsMessage
		if err := json.Unmarshal(msgBytes, &msg); err != nil {
			log.Printf("failed to parse websocket message: %v", err)
		}

		switch msg.Command {
		case "start_pomodoro":
			if err := s.session.RecordPomodoro(msg.Subject, io.Discard); err != nil {
				fmt.Fprintf(w, "failed to start pomodoro session for %q: %v", msg.Subject, err)
			}
		case "record_manual":
			if err := s.session.RecordManual(msg.Subject, msg.Hours); err != nil {
				fmt.Fprintf(w, "failed to record hours for %q: %v", msg.Subject, err)
			}
		default:
			fmt.Fprintln(w, "invalid command")
		}
	}
}

func (s *StudyServer) processGetRequest(w http.ResponseWriter, subject string) {
	hours, err := s.store.GetHours(subject)
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

	err = s.store.RecordHour(subject, h)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}
