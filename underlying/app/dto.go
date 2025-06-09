package app

import "time"

type QuizSettingsDto struct {
	Email         string
	QuestionCount int
	Applied       time.Time
	Quizzes       []QuizDto
	QuizChoice    QuizDto
}

type SummaryDto struct {
	Participants []ParticipantView
	QuizChoice   QuizDto
}

type QuizDto struct {
	Id          string
	Header      string
	Description string
}
