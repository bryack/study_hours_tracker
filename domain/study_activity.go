package domain

type StudyActivity struct {
	Subject string `json:"subject"`
	Hours   int    `json:"hours"`
}

type Report []StudyActivity
