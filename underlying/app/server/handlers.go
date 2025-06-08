package server

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/a-h/templ"
	"github.com/golangcollege/sessions"
	"github.com/google/uuid"
	"github.com/nefarius/cornelian/underlying/app"
	"github.com/nefarius/cornelian/underlying/app/access"
	"github.com/nefarius/cornelian/underlying/app/conf"
	"github.com/nefarius/cornelian/underlying/app/views"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Utility function to check if the user is logged in
func checkLoggedIn(session *sessions.Session, w http.ResponseWriter, r *http.Request) (string, bool) {
	email := session.GetString(r, "email")
	if email == "" {
		http.Error(w, "not logged in", http.StatusUnauthorized)
		return "", false
	}
	return email, true
}

func currentQuizPanelPage(session *sessions.Session, accessModule *access.CornelianModule) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		email, ok := checkLoggedIn(session, w, r)
		if !ok {
			return
		}
		currentQuiz := accessModule.QuestionService.GetQuiz()
		questions := accessModule.QuestionService.AllForAuthorInOpenStatus(email)

		accessModule.QuestionService.CreateIfNotExist(currentQuiz.Id, email)

		if len(questions) > 0 {
			var currentQuestion = questions[0]
			accessModule.QuestionService.StartAnswering(currentQuiz.Id, email, currentQuestion.ID)
			templ.Handler(views.CurrentQuizPanelPage(email, currentQuestion)).ServeHTTP(w, r)
		} else {
			var quizInfo = accessModule.QuestionService.QuizInfoRepository.GetQuizByIdAndEmail(currentQuiz.Id, email)
			templ.Handler(views.QuizFinishedPanelPage(quizInfo)).ServeHTTP(w, r)
		}

	}
}

func resetAnswersHandler(session *sessions.Session, accessModule *access.CornelianModule) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, ok := checkLoggedIn(session, w, r)
		if !ok {
			return
		}

		personId := r.URL.Query().Get("id")

		err := accessModule.QuestionService.ResetCurrentQuiz(personId)
		if err != nil {
			http.Error(w, "error saving quiz", http.StatusInternalServerError)
			return
		}

		participants := accessModule.QuestionService.GetParticipants()
		templ.Handler(views.ParticipantsPanelPage(participants)).ServeHTTP(w, r)
	}
}

func answerCurrentQuestionHandler(session *sessions.Session, accessModule *access.CornelianModule) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		email, ok := checkLoggedIn(session, w, r)
		if !ok {
			return
		}

		questionID := r.URL.Query().Get("id")

		answerText := r.FormValue("answertext")

		selected := r.FormValue("answer")
		if selected != "" {
			answerText = selected
		}

		currentQuiz := accessModule.QuestionService.GetQuiz()
		accessModule.QuestionService.HandleAnswer(currentQuiz.Id, email, questionID, answerText)

		err := accessModule.QuestionService.SaveAnswer(questionID, answerText, email)
		if err != nil {
			http.Error(w, "error saving answer", http.StatusInternalServerError)
			return
		}

		questions := accessModule.QuestionService.AllForAuthorInOpenStatus(email)
		print("active" + currentQuiz.Header)
		print("len " + fmt.Sprint(len(questions)))
		if len(questions) > 0 {
			var currentQuestion = questions[0]
			accessModule.QuestionService.StartAnswering(currentQuiz.Id, email, currentQuestion.ID)
			templ.Handler(views.CurrentQuizPanelPage(email, currentQuestion)).ServeHTTP(w, r)
		} else {
			var quizInfo = accessModule.QuestionService.QuizInfoRepository.GetQuizByIdAndEmail(currentQuiz.Id, email)
			templ.Handler(views.QuizFinishedPanelPage(quizInfo)).ServeHTTP(w, r)
		}
	}
}

func editQuestionHandler(session *sessions.Session, accessModule *access.CornelianModule) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		email, ok := checkLoggedIn(session, w, r)
		if !ok {
			return
		}

		questionID := r.FormValue("id")
		questionText := r.FormValue("questiontext")

		err := accessModule.QuestionService.UpdateQuestion(questionID, questionText, email)
		if err != nil {
			http.Error(w, "error saving question", http.StatusInternalServerError)
			return
		}

		indexPage(session, accessModule)(w, r)
	}
}

func updateSettingsHandler(session *sessions.Session, accessModule *access.CornelianModule) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		email, ok := checkLoggedIn(session, w, r)
		if !ok {
			return
		}

		questionCount, error := strconv.Atoi(r.FormValue("questioncount"))
		if error != nil {
			http.Error(w, "error parsing question count", http.StatusInternalServerError)
			return
		}

		var quizSettings = app.QuizSettings{
			QuestionCount: questionCount,
			Id:            primitive.NewObjectID(),
			Applied:       time.Now(),
			Current:       true,
			Email:         email,
		}

		err := accessModule.SettingsRepository.InsertSettings(quizSettings)
		if err != nil {
			http.Error(w, "error saving quiz", http.StatusInternalServerError)
			return
		}

		var settings = accessModule.SettingsRepository.GetAll()
		var currentSettings = filterSettingsByCurrentAndGetFirst(settings)

		templ.Handler(views.SettingsPage(currentSettings)).ServeHTTP(w, r)
	}
}

func filterSettingsByCurrentAndGetFirst(settings []app.QuizSettings) app.QuizSettings {
	var filtered []app.QuizSettings // Initialize an empty slice to hold the results

	for _, p := range settings { // Iterate over each product in the input slice
		if p.Current == true { // Check the boolean condition
			filtered = append(filtered, p) // If condition met, add to the filtered slice
		}
	}
	return filtered[0] // Return the new slice
}

func editQuizHandler(session *sessions.Session, accessModule *access.CornelianModule) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		email, ok := checkLoggedIn(session, w, r)
		if !ok {
			return
		}

		quizID := r.URL.Query().Get("id")
		quizText := r.FormValue("quizdescription")
		quizHeader := r.FormValue("quizheader")
		current := r.FormValue("quizcurrent")

		currentQuiz := current == "on"

		err := accessModule.QuestionService.UpdateQuiz(quizID, quizHeader, quizText, currentQuiz)
		if err != nil {
			http.Error(w, "error saving quiz", http.StatusInternalServerError)
			return
		}

		quizzes := accessModule.QuestionService.GetQuizzes()
		templ.Handler(views.QuizzesPanelPage(email, quizzes)).ServeHTTP(w, r)
	}
}

/**
 * We should put it into util method
 */
func createQuestion(email string, text string, answer1 string, answer2 string, answer3 string, answer4 string) app.Question {
	negRand := -rand.Intn(120)
	return app.Question{
		ID:        uuid.NewString(),
		From:      email,
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

func addNewQuizHandler(session *sessions.Session, accessModule *access.CornelianModule) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		email, ok := checkLoggedIn(session, w, r)
		if !ok {
			return
		}

		quizText := r.FormValue("quizdescription")
		quizHeader := r.FormValue("quizheader")

		// Iterate through potential question numbers.
		// You might need to know the max number of rows, or check for existence.
		// A simple loop for a fixed number of rows (e.g., 10) is common for fixed forms.
		// For dynamic forms, you might send a hidden field with the count of rows,
		// or look for patterns like "question_text_1", "question_text_2", etc.
		// 	var listOfQuestions = make([]app.Question, 0)
		var listOfQuestions = make([]app.Question, 0)
		for i := 0; i < len(listOfQuestions)+1; i++ {
			questionFieldName := fmt.Sprintf("question_%d", i)
			answerFieldName1 := fmt.Sprintf("answer1_%d", i)
			answerFieldName2 := fmt.Sprintf("answer2_%d", i)
			answerFieldName3 := fmt.Sprintf("answer3_%d", i)
			answerFieldName4 := fmt.Sprintf("answer4_%d", i)

			questionText := r.FormValue(questionFieldName)
			correctAnswer := r.FormValue(answerFieldName1)
			answer2 := r.FormValue(answerFieldName2)
			answer3 := r.FormValue(answerFieldName3)
			answer4 := r.FormValue(answerFieldName4)
			// If both fields for a given row are empty, assume no more rows
			if questionText == "" && correctAnswer == "" {
				break
			}

			// If only one is present, you might want to handle it as an error or partial input
			if questionText != "" || correctAnswer != "" {
				listOfQuestions = append(listOfQuestions, createQuestion(email, questionText, correctAnswer, answer2, answer3, answer4))
			}
		}
		// questions []app.Question
		err := accessModule.QuestionService.InsertQuizWithQuestionsAndMakeCurrent(quizHeader, quizText, email, listOfQuestions)
		if err != nil {
			http.Error(w, "error saving quiz", http.StatusInternalServerError)
			return
		}

		accessModule.CacheConf.Cache.Del(conf.CURRENT_TAG)

		quizzes := accessModule.QuestionService.GetQuizzes()
		templ.Handler(views.QuizzesPanelPage(email, quizzes)).ServeHTTP(w, r)
	}
}

func saveQuestionHandler(session *sessions.Session, accessModule *access.CornelianModule) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		email, ok := checkLoggedIn(session, w, r)
		if !ok {
			return
		}

		questionText := r.FormValue("questiontext")

		err := accessModule.QuestionService.AddQuestion(questionText, email)
		if err != nil {
			http.Error(w, "error saving question", http.StatusInternalServerError)
			return
		}

		indexPage(session, accessModule)(w, r)
	}
}

func answerQuestionHandler(session *sessions.Session, accessModule *access.CornelianModule) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		email, ok := checkLoggedIn(session, w, r)
		if !ok {
			return
		}

		questionID := r.URL.Query().Get("id")
		answerText := r.FormValue("answertext")

		err := accessModule.QuestionService.SaveAnswer(questionID, answerText, email)
		if err != nil {
			http.Error(w, "error saving answer", http.StatusInternalServerError)
			return
		}

		indexPage(session, accessModule)(w, r)
	}
}

func quizzesPanelPage(session *sessions.Session, accessModule *access.CornelianModule) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		email, ok := checkLoggedIn(session, w, r)
		if !ok {
			return
		}

		quizzes := accessModule.QuestionService.GetQuizzes()
		templ.Handler(views.QuizzesPanelPage(email, quizzes)).ServeHTTP(w, r)
	}
}

func participantsPanelPage(session *sessions.Session, accessModule *access.CornelianModule) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		_, ok := checkLoggedIn(session, w, r)
		if !ok {
			return
		}

		participants := accessModule.QuestionService.GetParticipants()
		templ.Handler(views.ParticipantsPanelPage(participants)).ServeHTTP(w, r)
	}
}

func answerQuestionPage(session *sessions.Session, accessModule *access.CornelianModule) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		email, ok := checkLoggedIn(session, w, r)
		if !ok {
			return
		}

		questionID := r.FormValue("id")
		question, err := accessModule.QuestionService.GetQuestion(questionID)
		if err != nil {
			http.Error(w, "question not found", http.StatusNotFound)
			return
		}
		templ.Handler(views.AnswerQuestion(email, question)).ServeHTTP(w, r)
	}
}

func addQuestionPage(session *sessions.Session, accessModule *access.CornelianModule) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		email, ok := checkLoggedIn(session, w, r)
		if !ok {
			return
		}

		quiz := accessModule.QuestionService.GetQuiz()
		templ.Handler(views.AddQuestion(email, quiz)).ServeHTTP(w, r)
	}
}

func addQuizPage(session *sessions.Session, accessModule *access.CornelianModule) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		email, ok := checkLoggedIn(session, w, r)
		if !ok {
			return
		}

		var settings = accessModule.SettingsRepository.GetAll()
		var currentSettings = filterSettingsByCurrentAndGetFirst(settings)
		numRows := currentSettings.QuestionCount
		indices := make([]int, numRows)
		for i := 0; i < numRows; i++ {
			indices[i] = i // Populate with actual index values
		}

		var quizCreationData = app.QuizCreationData{
			Email:           email,
			QuestionIndices: indices,
		}
		templ.Handler(views.AddQuizPage(quizCreationData)).ServeHTTP(w, r)
	}
}

func settingsPage(session *sessions.Session, accessModule *access.CornelianModule) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		_, ok := checkLoggedIn(session, w, r)
		if !ok {
			return
		}

		var settings = accessModule.SettingsRepository.GetAll()
		var currentSettings = filterSettingsByCurrentAndGetFirst(settings)

		templ.Handler(views.SettingsPage(currentSettings)).ServeHTTP(w, r)
	}
}

func editQuestionPage(session *sessions.Session, accessModule *access.CornelianModule) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		email, ok := checkLoggedIn(session, w, r)
		if !ok {
			return
		}

		questionID := r.URL.Query().Get("id")
		question, err := accessModule.QuestionService.GetQuestion(questionID)
		if err != nil {
			http.Error(w, "question not found", http.StatusNotFound)
			return
		}
		templ.Handler(views.EditQuestion(email, question)).ServeHTTP(w, r)
	}
}

func editQuizPage(session *sessions.Session, accessModule *access.CornelianModule) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		email, ok := checkLoggedIn(session, w, r)
		if !ok {
			return
		}

		quizID := r.URL.Query().Get("id")
		quiz, err := accessModule.QuestionService.GetQuizById(quizID)
		if err != nil {
			http.Error(w, "quiz not found", http.StatusNotFound)
			return
		}
		templ.Handler(views.EditQuiz(email, quiz)).ServeHTTP(w, r)
	}
}

func myQuestionsHandler(session *sessions.Session, accessModule *access.CornelianModule) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		email, ok := checkLoggedIn(session, w, r)
		if !ok {
			return
		}

		questions := accessModule.QuestionService.AllForAssignedTo(email)
		session.Put(r, "view", "mine")
		templ.Handler(views.Questions(questions)).ServeHTTP(w, r)
	}
}

func allQuestionsHandler(session *sessions.Session, accessModule *access.CornelianModule) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		_, ok := checkLoggedIn(session, w, r)
		if !ok {
			return
		}

		questions := accessModule.QuestionService.AllQuestions()
		session.Put(r, "view", "all")
		templ.Handler(views.Questions(questions)).ServeHTTP(w, r)
	}
}

func countOwnHandler(session *sessions.Session, accessModule *access.CornelianModule) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		email, ok := checkLoggedIn(session, w, r)
		if !ok {
			return
		}

		all := len(accessModule.QuestionService.AllForAuthorInOpenStatus(email))
		_, _ = w.Write([]byte(" (" + strconv.Itoa(all) + ")"))
	}
}

func countAllHandler(accessModule *access.CornelianModule) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		all := len(accessModule.QuestionService.CountQuestions())
		_, _ = w.Write([]byte(" (" + strconv.Itoa(all) + ")"))
	}
}
