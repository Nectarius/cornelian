package store

import (
	"math/rand"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/nefarius/cornelian/underlying/app"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetDefaultQuizData() app.Quiz {

	var listOfQuestions = make([]app.Question, 0)
	talk1 := app.Talk{
		ID:      "1",
		Title:   "Templ",
		Authors: []string{"aeneole@gmail.com"},
	}

	talks := []app.Talk{talk1}
	gofakeit.Seed(time.Now().UnixMilli())

	for j := 0; j < 15; j++ {
		rnd := rand.Intn(len(talks))
		negRand := -rand.Intn(120)
		q := app.Question{
			ID:        uuid.NewString(),
			Talk:      talks[rnd],
			From:      gofakeit.Name(),
			Text:      gofakeit.Question(),
			CreatedAt: time.Now().Add(time.Minute * time.Duration(negRand)),
			Status:    app.StatusOpen,
		}
		listOfQuestions = append(listOfQuestions, q)
	}

	return app.Quiz{
		Id:          primitive.NewObjectID(),
		Header:      "Quiz",
		Description: "Quiz",
		Tag:         "Test",
		Questions:   listOfQuestions,
	}
}
