package server

import (
	"net/http"

	"github.com/golangcollege/sessions"
	"github.com/markbates/goth/gothic"
	"github.com/nefarius/cornelian/underlying/app"
	"github.com/nefarius/cornelian/underlying/app/access"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func authCallbackHandler(session *sessions.Session, accessModule *access.CornelianModule) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		user, err := gothic.CompleteUserAuth(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		var personRepository = accessModule.PersonRepository

		personRepository.InsertPersonIfNotPresent(app.Person{
			Id:    primitive.NewObjectID(),
			Email: user.Email,
		})

		var questionService = accessModule.QuestionService

		questionService.AssignQuizIfApplicable(user.Email)

		session.Put(r, "email", user.Email)
		session.Put(r, "view", "all")
		http.Redirect(w, r, "/", 302)
	}
}

func authStartHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		gothic.BeginAuthHandler(w, r)
	}
}
