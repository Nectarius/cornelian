package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/a-h/templ"
	"github.com/golangcollege/sessions"
	"github.com/nefarius/cornelian/underlying/app"
	"github.com/nefarius/cornelian/underlying/app/access"
	"github.com/nefarius/cornelian/underlying/app/views"
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

func answerCurrentQuestionHandler(session *sessions.Session, accessModule *access.CornelianModule) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		email, ok := checkLoggedIn(session, w, r)
		if !ok {
			return
		}

		questionID := r.URL.Query().Get("id")
		answerText := r.FormValue("answertext")
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

func addNewQuizHandler(session *sessions.Session, accessModule *access.CornelianModule) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		email, ok := checkLoggedIn(session, w, r)
		if !ok {
			return
		}

		quizText := r.FormValue("quizdescription")
		quizHeader := r.FormValue("quizheader")

		err := accessModule.QuestionService.InsertQuizAndMakeCurrent(quizHeader, quizText, email)
		if err != nil {
			http.Error(w, "error saving quiz", http.StatusInternalServerError)
			return
		}

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

		templ.Handler(views.AddQuizPage(email)).ServeHTTP(w, r)
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
		all := len(accessModule.QuestionService.AllInStatus(app.StatusOpen))
		_, _ = w.Write([]byte(" (" + strconv.Itoa(all) + ")"))
	}
}
