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

func editQuestionHandler(session *sessions.Session, accessModule *access.CornelianModule) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		email := session.GetString(r, "email")
		if email == "" {
			http.Error(w, "not logged in", 401)
			return
		}

		var questionId = r.FormValue("id")
		questionText := r.FormValue("questiontext")

		var questionService = accessModule.QuestionService

		err := questionService.UpdateQuestion(questionId, questionText, email)

		if err != nil {
			http.Error(w, "error saving answer", http.StatusInternalServerError)
			return
		}

		indexPage(session, accessModule)(w, r)
	}
}

func editQuizHandler(session *sessions.Session, accessModule *access.CornelianModule) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		email := session.GetString(r, "email")
		if email == "" {
			http.Error(w, "not logged in", 401)
			return
		}

		var quizId = r.URL.Query().Get("id")

		quizText := r.FormValue("quizdescription")
		quizHeader := r.FormValue("quizheader")

		current := r.FormValue("quizcurrent")

		var currentQuiz bool
		if current == "on" {
			currentQuiz = true
		} else {
			currentQuiz = false
		}

		fmt.Println("current + " + current)

		var questionService = accessModule.QuestionService

		err := questionService.UpdateQuiz(quizId, quizHeader, quizText, currentQuiz)

		if err != nil {
			http.Error(w, "error saving answer", http.StatusInternalServerError)
			return
		}

		quizzes := questionService.GetQuizzes()

		templ.Handler(views.QuizzesPanelPage(email, quizzes)).ServeHTTP(w, r)
	}
}

func addNewQuizHandler(session *sessions.Session, accessModule *access.CornelianModule) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		email := session.GetString(r, "email")
		if email == "" {
			http.Error(w, "not logged in", 401)
			return
		}

		//questionID := r.FormValue("id")
		quizText := r.FormValue("quizdescription")
		quizHeader := r.FormValue("quizheader")

		var questionService = accessModule.QuestionService

		err := questionService.InsertQuizAndMakeCurrent(quizHeader, quizText, email)

		if err != nil {
			http.Error(w, "error saving answer", http.StatusInternalServerError)
			return
		}

		quizzes := questionService.GetQuizzes()

		templ.Handler(views.QuizzesPanelPage(email, quizzes)).ServeHTTP(w, r)
	}
}

func saveQuestionHandler(session *sessions.Session, accessModule *access.CornelianModule) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		email := session.GetString(r, "email")
		if email == "" {
			http.Error(w, "not logged in", 401)
			return
		}

		//questionID := r.FormValue("id")
		questionText := r.FormValue("questiontext")

		var questionService = accessModule.QuestionService

		err := questionService.AddQuestion(questionText, email)

		if err != nil {
			http.Error(w, "error saving answer", http.StatusInternalServerError)
			return
		}

		indexPage(session, accessModule)(w, r)
	}
}

func answerQuestionHandler(session *sessions.Session, accessModule *access.CornelianModule) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		email := session.GetString(r, "email")
		if email == "" {
			http.Error(w, "not logged in", 401)
			return
		}

		var questionId = r.URL.Query().Get("id")
		answerText := r.FormValue("answertext")

		var questionService = accessModule.QuestionService

		err := questionService.SaveAnswer(questionId, answerText, email)

		if err != nil {
			http.Error(w, "error saving answer", http.StatusInternalServerError)
			return
		}

		indexPage(session, accessModule)(w, r)
	}
}

func quizzesPanelPage(session *sessions.Session, accessModule *access.CornelianModule) func(w http.ResponseWriter, r *http.Request) {
	var questionService = accessModule.QuestionService
	return func(w http.ResponseWriter, r *http.Request) {
		email := session.GetString(r, "email")
		// check permissions
		if email == "" {
			http.Error(w, "not logged in", 401)
			return
		}
		// questionID := r.FormValue("id")
		quizzes := questionService.GetQuizzes()

		templ.Handler(views.QuizzesPanelPage(email, quizzes)).ServeHTTP(w, r)
	}
}

func answerQuestionPage(session *sessions.Session, accessModule *access.CornelianModule) func(w http.ResponseWriter, r *http.Request) {
	var questionRepository = accessModule.QuestionService
	return func(w http.ResponseWriter, r *http.Request) {
		email := session.GetString(r, "email")
		if email == "" {
			http.Error(w, "not logged in", 401)
			return
		}
		questionID := r.FormValue("id")
		question, err := questionRepository.GetQuestion(questionID)
		if err != nil {
			http.Error(w, "question not found", 404)
			return
		}
		templ.Handler(views.AnswerQuestion(email, question)).ServeHTTP(w, r)
	}
}

func addQuestionPage(session *sessions.Session, accessModule *access.CornelianModule) func(w http.ResponseWriter, r *http.Request) {
	var questionService = accessModule.QuestionService
	return func(w http.ResponseWriter, r *http.Request) {
		email := session.GetString(r, "email")
		if email == "" {
			http.Error(w, "not logged in", 401)
			return
		}
		//questionID := r.FormValue("id")
		var quiz = questionService.GetQuiz()

		templ.Handler(views.AddQuestion(email, quiz)).ServeHTTP(w, r)
	}
}

func addQuizPage(session *sessions.Session, accessModule *access.CornelianModule) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		email := session.GetString(r, "email")
		if email == "" {
			http.Error(w, "not logged in", 401)
			return
		}

		templ.Handler(views.AddQuizPage(email)).ServeHTTP(w, r)
	}
}

func editQuestionPage(session *sessions.Session, accessModule *access.CornelianModule) func(w http.ResponseWriter, r *http.Request) {
	var questionService = accessModule.QuestionService
	return func(w http.ResponseWriter, r *http.Request) {
		email := session.GetString(r, "email")
		if email == "" {
			http.Error(w, "not logged in", 401)
			return
		}
		var questionId = r.URL.Query().Get("id")
		var question, err = questionService.GetQuestion(questionId)
		if err != nil {
			http.Error(w, "question not found", 404)
			return
		}
		templ.Handler(views.EditQuestion(email, question)).ServeHTTP(w, r)
	}
}

func editQuizPage(session *sessions.Session, accessModule *access.CornelianModule) func(w http.ResponseWriter, r *http.Request) {
	var questionService = accessModule.QuestionService
	return func(w http.ResponseWriter, r *http.Request) {
		email := session.GetString(r, "email")
		if email == "" {
			http.Error(w, "not logged in", 401)
			return
		}
		var quizId = r.URL.Query().Get("id")
		var quiz, err = questionService.GetQuizById(quizId)
		if err != nil {
			http.Error(w, "question not found", 404)
			return
		}
		templ.Handler(views.EditQuiz(email, quiz)).ServeHTTP(w, r)
	}
}

func myQuestionsHandler(session *sessions.Session, accessModule *access.CornelianModule) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		email := session.GetString(r, "email")
		if email == "" {
			http.Error(w, "not logged in", 401)
			return
		}
		var repository = accessModule.QuestionService
		questions := repository.AllForAssignedTo(email)
		session.Put(r, "view", "mine")
		templ.Handler(views.Questions(questions)).ServeHTTP(w, r)
	}
}

func allQuestionsHandler(session *sessions.Session, accessModule *access.CornelianModule) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		email := session.GetString(r, "email")
		if email == "" {
			http.Error(w, "not logged in", 401)
			return
		}
		var repository = accessModule.QuestionService
		questions := repository.AllQuestions()
		session.Put(r, "view", "all")
		templ.Handler(views.Questions(questions)).ServeHTTP(w, r)
	}
}

func countOwnHandler(session *sessions.Session, accessModule *access.CornelianModule) func(w http.ResponseWriter, r *http.Request) {
	var questionService = accessModule.QuestionService
	questionService.AllQuestions()
	return func(w http.ResponseWriter, r *http.Request) {
		email := session.GetString(r, "email")
		all := len(questionService.AllForAuthorInOpenStatus(email))
		_, _ = w.Write([]byte(" (" + strconv.Itoa(all) + ")"))
	}
}

func countAllHandler(accessModule *access.CornelianModule) func(w http.ResponseWriter, r *http.Request) {
	var repository = accessModule.QuestionService
	return func(w http.ResponseWriter, r *http.Request) {
		all := len(repository.AllInStatus(app.StatusOpen))
		_, _ = w.Write([]byte(" (" + strconv.Itoa(all) + ")"))
	}
}

// func deleteQuestionHandler(session *sessions.Session, db *store.InMem) func(w http.ResponseWriter, r *http.Request) {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		email := session.GetString(r, "email")
// 		if email == "" {
// 			http.Error(w, "not logged in", 401)
// 			return
// 		}

// 		questionID := r.URL.Query().Get("id")
// 		db.Delete(questionID)
// 		if session.GetString(r, "view") == "mine" {
// 			templ.Handler(views.Questions(db.AllForAssignedTo(email))).ServeHTTP(w, r)
// 		} else {
// 			templ.Handler(views.Questions(db.All())).ServeHTTP(w, r)
// 		}

// 	}
// }
