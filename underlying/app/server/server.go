package server

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/golangcollege/sessions"
	"github.com/nefarius/cornelian/underlying/app/access"
	"github.com/nefarius/cornelian/underlying/app/store"
	"github.com/nefarius/cornelian/underlying/app/views"
)

func StartServer(session *sessions.Session, db *store.InMem, accessModule *access.CornelianModule) {
	// Set-up chi router with middleware
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(session.Enable)

	// Page specific handlers
	r.Get("/", indexPage(session, accessModule))
	r.Get("/login", templ.Handler(views.Login()).ServeHTTP)
	r.Get("/answer", answerQuestionPage(session, accessModule))
	r.Get("/add-question", addQuestionPage(session, accessModule))
	r.Get("/edit-question", editQuestionPage(session, accessModule))

	// Login handlers
	r.Get("/auth", authStartHandler())
	r.Get("/auth/{provider}/callback", authCallbackHandler(session))
	r.Get("/logout", logoutHandler(session))

	// API handlers
	r.Get("/countall", countAllHandler(accessModule))
	r.Get("/countmine", countOwnHandler(session, accessModule))

	r.Get("/all", allQuestionsHandler(session, accessModule))
	r.Get("/mine", myQuestionsHandler(session, accessModule))

	r.Post("/save-question", saveQuestionHandler(session, accessModule))
	r.Post("/update-question", editQuestionHandler(session, accessModule))
	r.Post("/answerquestion", answerQuestionHandler(session, accessModule))
	// r.Delete("/delete", deleteQuestionHandler(session, db))

	// Start plain HTTP listener
	_ = http.ListenAndServe(":3000", r)
}

func indexPage(session *sessions.Session, accessModule *access.CornelianModule) func(w http.ResponseWriter, r *http.Request) {
	var questionService = accessModule.QuestionService
	return func(w http.ResponseWriter, r *http.Request) {
		email := session.GetString(r, "email")
		if email != "" {
			if session.GetString(r, "view") == "mine" {
				templ.Handler(views.Index(email, questionService.AllForAssignedTo(email))).ServeHTTP(w, r)
			} else {
				templ.Handler(views.Index(email, questionService.AllQuestions())).ServeHTTP(w, r)
			}
			return
		}
		templ.Handler(views.Index("", nil)).ServeHTTP(w, r)
	}
}
