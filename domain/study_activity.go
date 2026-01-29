package domain

import "errors"

var ErrSubjectNotFound = errors.New("subject not found")

type StudyActivity struct {
	Subject string `json:"subject"`
	Hours   int    `json:"hours"`
}

type Report []StudyActivity
