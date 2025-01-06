package store

import (
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/nefarius/cornelian/underlying/app"
	"github.com/nefarius/cornelian/underlying/app/conf"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetDefaultQuizData() app.Quiz {

	var listOfQuestions = make([]app.Question, 0)

	negRand := -rand.Intn(120)
	q1 := app.Question{
		ID:        uuid.NewString(),
		From:      "nefarius",
		Text:      "Команда становившаяся в период с 2003 - 2008 Наибольшее число раз чемпионом. <br /> ЦСКА Локомотив Рубин Зенит",
		CreatedAt: time.Now().Add(time.Minute * time.Duration(negRand)),
		Status:    app.StatusOpen,
	}

	listOfQuestions = append(listOfQuestions, q1)

	q2 := app.Question{
		ID:        uuid.NewString(),
		From:      "nefarius",
		Text:      " В какой команде РПЛ начинал играть в России Мигель Данни  - Динами Зенит ЦСКА Локомотив",
		CreatedAt: time.Now().Add(time.Minute * time.Duration(negRand)),
		Status:    app.StatusOpen,
	}

	listOfQuestions = append(listOfQuestions, q2)
	//}

	return app.Quiz{
		Id:          primitive.NewObjectID(),
		Header:      "РФПЛ 2000 - 2024",
		Description: "Футбольные вопроссы",
		Tag:         conf.CURRENT_TAG,
		Questions:   listOfQuestions,
	}
}
