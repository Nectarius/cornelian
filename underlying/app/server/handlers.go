package server

import (
	"net/http"
	"strconv"

	"github.com/a-h/templ"
	"github.com/golangcollege/sessions"
	"github.com/nefarius/cornelian/underlying/app"
	"github.com/nefarius/cornelian/underlying/app/access"
	"github.com/nefarius/cornelian/underlying/app/views"
)

func answerQuestionHandler(session *sessions.Session, accessModule *access.CornelianModule) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		email := session.GetString(r, "email")
		if email == "" {
			http.Error(w, "not logged in", 401)
			return
		}

		questionID := r.FormValue("id")
		answerText := r.FormValue("answertext")

		var questionService = accessModule.QuestionService

		err := questionService.SaveAnswer(questionID, answerText, email)

		if err != nil {
			http.Error(w, "error saving answer", http.StatusInternalServerError)
			return
		}

		indexPage(session, accessModule)(w, r)
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
	var repository = accessModule.QuestionService
	return func(w http.ResponseWriter, r *http.Request) {
		email := session.GetString(r, "email")
		all := len(repository.AllForAuthorInStatus(email, app.StatusOpen))
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
