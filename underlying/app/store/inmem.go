package store

import (
	"fmt"
	"log/slog"
	"slices"
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/nefarius/cornelian/underlying/app"
)

type InMem struct {
	db []app.Question
}

func NewInMem() *InMem {
	db := make([]app.Question, 0)
	return &InMem{db: db}
}

func (i *InMem) Create(q app.Question) {
	i.db = append(i.db, q)
}
func (i *InMem) All() []app.Question {
	out := i.db
	sort.Slice(out, func(i, j int) bool {
		return out[i].CreatedAt.After(out[j].CreatedAt)
		//return out[i].CreatedAt.After(out[j].CreatedAt) && out[i].Status == app.StatusOpen
	})
	return out
}

func (i *InMem) AllInStatus(status app.Status) []app.Question {
	out := make([]app.Question, 0)
	for _, q := range i.db {
		if q.Status == status {
			out = append(out, q)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		//return (out[j].Status == app.StatusOpen && out[i].Status != app.StatusOpen) && out[i].CreatedAt.After(out[j].CreatedAt)
		return out[i].CreatedAt.After(out[j].CreatedAt)
	})
	return out
}

func (i *InMem) Delete(id string) {
	slog.Info("delete question by id", slog.String("id", id))
	for j := range i.db {
		if i.db[j].ID == id {
			i.db = slices.Delete(i.db, j, j+1)
			break
		}
	}
}

func (i *InMem) FillWithData(q []app.Question) {
	for _, question := range q {
		i.db = append(i.db, question) // note the = instead of :=
	}
}

func (i *InMem) Get(id string) (app.Question, error) {
	for _, q := range i.db {
		if q.ID == id {
			return q, nil
		}
	}
	return app.Question{}, fmt.Errorf("question identified by %v not found", id)
}

func (i *InMem) SaveAnswer(id string, text string, answeredBy string) error {
	q, err := i.Get(id)
	if err != nil {
		return err
	}
	q.Answers = append(q.Answers, app.Answer{ID: uuid.NewString(), AnsweredBy: answeredBy, Text: text, AnsweredAt: time.Now()})
	q.Status = app.StatusAnswered

	// Delete the old one and add the new
	i.Delete(id)
	i.Create(q)
	return nil
}
