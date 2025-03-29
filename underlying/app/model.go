package app

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Status string

const (
	StatusOpen     = Status("open")
	StatusAnswered = Status("answered")
	StatusRemoved  = Status("removed")
)

type Person struct {
	Id    primitive.ObjectID
	Email string
	Admin bool
}

type Question struct {
	ID        string
	From      string
	Text      string
	CreatedAt time.Time
	Status    Status
	Answers   []Answer
}

type Answer struct {
	ID         string
	Text       string
	AnsweredBy string
	AnsweredAt time.Time
}

type AnswerInfo struct {
	ID         string
	QuestionId string
	Text       string
	Current    bool
	Started    time.Time
	Completed  time.Time
}

type Quiz struct {
	Id          primitive.ObjectID
	Header      string
	Description string
	Active      bool
	Current     bool
	Tag         string
	Creator     string
	Questions   []Question
	AssignedTo  []string
}

type QuizInfo struct {
	Id        primitive.ObjectID
	QuizId    primitive.ObjectID
	Email     string
	Started   time.Time
	Completed time.Time
	Answers   []AnswerInfo
}
