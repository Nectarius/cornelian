package store

import (
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/nefarius/cornelian/underlying/app"
	"github.com/nefarius/cornelian/underlying/app/conf"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func createQuestion(text string, answer1 string, answer2 string, answer3 string, answer4 string) app.Question {
	negRand := -rand.Intn(120)
	return app.Question{
		ID:        uuid.NewString(),
		From:      "nefarius",
		Text:      text,
		CreatedAt: time.Now().Add(time.Minute * time.Duration(negRand)),
		Status:    app.StatusOpen,
		Answers:   make([]app.Answer, 0),
		AnswerChoices: []app.AnswerChoice{
			{ID: uuid.NewString(), Text: answer1, CorrectResponse: true},
			{ID: uuid.NewString(), Text: answer2, CorrectResponse: false},
			{ID: uuid.NewString(), Text: answer3, CorrectResponse: false},
			{ID: uuid.NewString(), Text: answer4, CorrectResponse: false},
		},
	}
}

func GetDefaultQuizData2() app.Quiz {

	var listOfQuestions = make([]app.Question, 0)

	listOfQuestions = append(listOfQuestions, createQuestion("Самый вместительный стадион Англии", "Уэмбли", "Олд Траффорд", "Энфилд", "Тоттенхэм Хотспур Стэдиум"))
	listOfQuestions = append(listOfQuestions, createQuestion("У кого из этих клубов наибольшее количество чемпионских сезонов подряд", "Манчестер Сити", "Манчестер Юнайтед", "Арсенал", "Ливерпуль"))
	listOfQuestions = append(listOfQuestions, createQuestion("Из какого города футбольный клуб Тоттенхем", "Лондон", "Ливерпуль", "Бирмингем", "Манчестер"))
	listOfQuestions = append(listOfQuestions, createQuestion("В каком клубе Англии играл Андрей Аршавин", "Арсенал", "Тоттенхем", "Манчестер Юнайтед", "Челси"))
	listOfQuestions = append(listOfQuestions, createQuestion("Лучший бомбардир в истории Арсенала", "Тьери Анри", "Робин Ван Перси", "Руд Ван Нистелрой", "Джон Рэдфорд"))

	listOfQuestions = append(listOfQuestions, createQuestion("Лучший бомбардир Премьер-лиги", "Алан Ширер", "Уэйн Руни", "Гарри Кейн", "Энди Коул"))
	listOfQuestions = append(listOfQuestions, createQuestion("Кто из этих клубов никогда не был чемпионом Англии", "Тоттенхем Хотспур", "Лестер Сити", "Блэкберн Роверс", "Ливерпуль"))
	listOfQuestions = append(listOfQuestions, createQuestion("В каком клубе играл Роман Павлюченко", "Тоттенхем", "Манчестер Юнайтед", "Арсенал", "Челси"))
	listOfQuestions = append(listOfQuestions, createQuestion("Какой из перечисленных Российский футболист играл за Эвертон", "Андрей  Канчельскис", "Роман Павлюченко", "Александр Мостовой", "Роман Зобнин"))
	listOfQuestions = append(listOfQuestions, createQuestion("Матч каких команд называют северолондонским дерби", "Тоттенхем Хотспур и Арсенала", "Челси и Арсенала", "Тоттенхем Хотспур и Челси", "Вест Хэм Юнайтед и Тоттенхем Хотспур"))

	listOfQuestions = append(listOfQuestions, createQuestion("Какой клуб Англии носит прозвище  «молотобойцы»", "Вест Хэм Юнайтед", "Лидс", "Эвертон", "Саутгемптон"))
	listOfQuestions = append(listOfQuestions, createQuestion("Кто становился наибольшее количество раз чемпионом в период с 2006 по 2011 год", "Манчестер Юнайтед", "Челси", "Арсенал", "Манчестер Сити"))
	listOfQuestions = append(listOfQuestions, createQuestion("У какого игрока Манчестер Юнайтед было прозвище Валлийский волшебник", "Райан Гиггз", "Эдвард Шерингем", "Пол Скоулз", "Гари Нэвилл"))
	listOfQuestions = append(listOfQuestions, createQuestion("Кто был капитаном Ливерпуля с 2003 по 2015 год", "Стивен Джеррард", "Хавьер Маскерано", "Джейми Каррагер", "Джордан Хендерсон"))
	listOfQuestions = append(listOfQuestions, createQuestion("Кто становился наибольшее количество раз чемпионом в период с 2000/01 по 2006/07 год", "Манчестер Юнайтед", "Челси", "Арсенал", "Манчестер Сити"))

	var assignedTo = make([]string, 0)
	assignedTo = append(assignedTo, "aeneole@gmail.com")
	return app.Quiz{
		Id:          primitive.NewObjectID(),
		Header:      "АПЛ 2000 - 2025",
		Description: "АПЛ. Футбольные вопроссы с вариантами ответов",
		Tag:         conf.CURRENT_TAG,
		Questions:   listOfQuestions,
		Active:      true,
		Current:     true,
		AssignedTo:  assignedTo,
	}
}

/*func GetDefaultQuizData() app.Quiz {

var listOfQuestions = make([]app.Question, 0)

listOfQuestions = append(listOfQuestions, createQuestion("Команда становившаяся в период с 2003 - 2008 Наибольшее число раз чемпионом - ЦСКА Локомотив Рубин Зенит"))
listOfQuestions = append(listOfQuestions, createQuestion("В какой команде РПЛ начинал играть в России Мигель Данни - Динами Зенит ЦСКА Локомотив"))
listOfQuestions = append(listOfQuestions, createQuestion("Сколько раз Рубин становился чемпионом ? 3, 0 , 2 , 1"))
listOfQuestions = append(listOfQuestions, createQuestion("Кто обыгрывал в матче за супер кубок Манчестер Юнайтед - ЦСКА Зенит Спартак Локомотив"))
listOfQuestions = append(listOfQuestions, createQuestion("Сейчас Луча́но Спалле́тти тренер сборной Италии,а когда то он тренировал клуб РПЛ. Сможете его вспомнить ? - Зенит Спартак Локомотив Динамо"))

listOfQuestions = append(listOfQuestions, createQuestion("Лучший бомбардир 2000 года Дмитрий Лоськов Дмитрий Кириченко Егор Титов Сергей Семак"))
listOfQuestions = append(listOfQuestions, createQuestion("Сколько раз Локомотив становился чемпионом ? 3, 4 , 2 , 1"))
listOfQuestions = append(listOfQuestions, createQuestion("Сейчас Мурат Якин тренер сборной Швейцарии,а когда то он тренировал клуб РПЛ. Сможете его вспомнить ? - Зенит Спартак Локомотив Динамо"))
listOfQuestions = append(listOfQuestions, createQuestion("Первая команда России с вышедшим на футбольном поле составом не только из легионеров но даже без игроков с постсоветского пространства Динамо Сатурн ЦСКА Рубин"))
listOfQuestions = append(listOfQuestions, createQuestion("Кто стал лучшим бомбардиром чемпионата России по футболу 2013 - 2014? Дзюба Мовсисян Кержаков Думбия"))

listOfQuestions = append(listOfQuestions, createQuestion(" У кого был самый молодой состав РПЛ в 2023 год? Краснодар Локомотив Спартак Крылья"))
listOfQuestions = append(listOfQuestions, createQuestion("Кто стал самым дорогим зимним новичком РПЛ зимой в 2024 ? Педро Угальде Кастаньо Артур"))
listOfQuestions = append(listOfQuestions, createQuestion("Зимой Премьер-Лигу пополнил футболист-однофамилец легендарного бразильца. Кто он? - Пеле Зико Роналдо Кака"))
listOfQuestions = append(listOfQuestions, createQuestion("Маршрут из «Зенита» в «Сочи» для РПЛ — дело привычное. Этой зимой по нему последовали двое. Кто остался в Петербурге? Сутормин. Чистяков. Ерохин Глушенков"))
listOfQuestions = append(listOfQuestions, createQuestion("Команда становившаяся в период с 2011 - 2016 Наибольшее число раз чемпионом - ЦСКА Локомотив Зенит Спартак"))

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
}*/
