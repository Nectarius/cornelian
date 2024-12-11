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

type Question struct {
	ID        string
	Talk      Talk
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

type Talk struct {
	ID      string
	Title   string
	Authors []string
}

type Quiz struct {
	Id          primitive.ObjectID
	Header      string
	Description string
	Tag         string
	Questions   []Question
}
