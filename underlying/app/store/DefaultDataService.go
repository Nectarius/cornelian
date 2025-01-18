package store

import (
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/nefarius/cornelian/underlying/app"
	"github.com/nefarius/cornelian/underlying/app/conf"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func createQuestion(text string) app.Question {
	negRand := -rand.Intn(120)
	return app.Question{
		ID:        uuid.NewString(),
		From:      "nefarius",
		Text:      text,
		CreatedAt: time.Now().Add(time.Minute * time.Duration(negRand)),
		Status:    app.StatusOpen,
		Answers:   make([]app.Answer, 0),
	}
}

func GetDefaultQuizData() app.Quiz {

	var listOfQuestions = make([]app.Question, 0)

	listOfQuestions = append(listOfQuestions, createQuestion("Команда становившаяся в период с 2003 - 2008 Наибольшее число раз чемпионом - ЦСКА Локомотив Рубин Зенит"))
	listOfQuestions = append(listOfQuestions, createQuestion("В какой команде РПЛ начинал играть в России Мигель Данни - Динами Зенит ЦСКА Локомотив"))
	listOfQuestions = append(listOfQuestions, createQuestion("Сколько раз Рубин становился чемпионом ? 3, 0 , 2 , 1"))
	listOfQuestions = append(listOfQuestions, createQuestion("Кто обыгрывал в матче за супер кубок Манчестер Юнайтед - ЦСКА Зенит Спартак Локомотив"))
	listOfQuestions = append(listOfQuestions, createQuestion("Сейчас Луча́но Спалле́тти тренер сборной Италии,а когда то он тренировал клуб РПЛ. Сможете его вспомнить ?"))

	listOfQuestions = append(listOfQuestions, createQuestion("Лучший бомбардир 2000 года Дмитрий Лоськов Дмитрий Кириченко Егор Титов Сергей Семак"))
	listOfQuestions = append(listOfQuestions, createQuestion("Сколько раз Локомотив становился чемпионом ? 3, 4 , 2 , 1"))
	listOfQuestions = append(listOfQuestions, createQuestion("Сейчас Мурат Якин тренер сборной Швейцарии,а когда то он тренировал клуб РПЛ. Сможете его вспомнить ?"))
	listOfQuestions = append(listOfQuestions, createQuestion("Первая команда России с вышедшим на футбольном поле составом не только из легионеров но даже без игроков с постсоветского пространства Динамо Сатурн ЦСКА Рубин"))
	listOfQuestions = append(listOfQuestions, createQuestion("Кто стал лучшим бомбардиром чемпионата России по футболу 2013 - 2014? Дзюба Мовсисян Кержаков Думбия"))

	listOfQuestions = append(listOfQuestions, createQuestion(" У кого был самый молодой состав РПЛ в 2023 год? Краснодар Локомотив Спартак Крылья"))
	listOfQuestions = append(listOfQuestions, createQuestion("Кто стал самым дорогим зимним новичком РПЛ зимой в 2024 ? Педро Угальде Кастаньо Артур"))
	listOfQuestions = append(listOfQuestions, createQuestion("Зимой Премьер-Лигу пополнил футболист-однофамилец легендарного бразильца. Кто он?"))
	listOfQuestions = append(listOfQuestions, createQuestion("Маршрут из «Зенита» в «Сочи» для РПЛ — дело привычное. Этой зимой по нему последовали двое. Кто остался в Петербурге?"))
	listOfQuestions = append(listOfQuestions, createQuestion("Команда становившаяся в период с 2011 - 2016 Наибольшее число раз чемпионом ЦСКА Локомотив Зенит Спартак"))

	var assignedTo = make([]string, 0)
	assignedTo = append(assignedTo, "aeneole@gmail.com")
	return app.Quiz{
		Id:          primitive.NewObjectID(),
		Header:      "РФПЛ 2000 - 2024",
		Description: "Футбольные вопроссы",
		Tag:         conf.CURRENT_TAG,
		Questions:   listOfQuestions,
		Active:      true,
		Current:     true,
		AssignedTo:  assignedTo,
	}
}
