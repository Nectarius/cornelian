package server

import (
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/golangcollege/sessions"
	"github.com/nefarius/cornelian/underlying/app/access"
	"github.com/nefarius/cornelian/underlying/app/store"
	"github.com/nefarius/cornelian/underlying/app/views"
	"golang.org/x/crypto/acme/autocert"
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
	r.Get("/add-quiz", addQuizPage(session, accessModule))

	// Login handlers
	r.Get("/auth", authStartHandler())
	r.Get("/auth/{provider}/callback", authCallbackHandler(session, accessModule))
	r.Get("/logout", logoutHandler(session))
	r.Get("/quizzes-panel", quizzesPanelPage(session, accessModule))
	r.Get("/edit-quiz", editQuizPage(session, accessModule))

	// API handlers
	r.Get("/countall", countAllHandler(accessModule))
	r.Get("/countmine", countOwnHandler(session, accessModule))

	r.Get("/all", allQuestionsHandler(session, accessModule))
	r.Get("/mine", myQuestionsHandler(session, accessModule))

	r.Post("/save-question", saveQuestionHandler(session, accessModule))
	r.Post("/add-new-quiz", addNewQuizHandler(session, accessModule))

	r.Post("/update-question", editQuestionHandler(session, accessModule))
	r.Post("/update-quiz", editQuizHandler(session, accessModule))

	r.Post("/answerquestion", answerQuestionHandler(session, accessModule))

	//httpsListenAndServeWithLetsEncrypt(r)
	httpsListenAndServe(r)
	//localListenAndServe(r)
}

func httpsListenAndServeWithLetsEncrypt(r *chi.Mux) {
	// localhost
	//certFile := "cert.pem"
	//keyFile := "key.pem"

	// kornelian host
	m := autocert.Manager{
		Prompt:     autocert.AcceptTOS,                      // Choose how to handle challenges (e.g., HTTP-01)
		HostPolicy: autocert.HostWhitelist("kornelian.com"), // Add your domain here
		Cache:      autocert.DirCache("certs"),              // Directory to store certificates
	}

	// Create a TLS configuration

	config := &tls.Config{GetCertificate: m.GetCertificate}
	server := &http.Server{
		Addr:      ":443", // Listen on port 443 (HTTPS)
		Handler:   r,
		TLSConfig: config,
	}

	// Start plain HTTP listener
	//_ = http.ListenAndServe(":5120", r)
	fmt.Println("Configured TLS with autocert.Manager...")
	fmt.Println("Server listening on HTTPS...")
	err := server.ListenAndServeTLS("", "") // Use cert and key from config
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}

func httpsListenAndServe(r *chi.Mux) {
	// localhost
	//certFile := "cert.pem"
	//keyFile := "key.pem"

	// kornelian host
	certFile := "kornelian.com.pem"
	keyFile := "kornelian.com.key"

	// Create a TLS configuration
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		//fmt.Println("Error loading certificate:", err)
		return
	}
	config := &tls.Config{Certificates: []tls.Certificate{cert}}
	server := &http.Server{
		Addr:      ":443", // Listen on port 443 (HTTPS)
		Handler:   r,
		TLSConfig: config,
	}

	// Start plain HTTP listener
	//_ = http.ListenAndServe(":5120", r)

	fmt.Println("Server listening on HTTPS...")
	err = server.ListenAndServeTLS("", "") // Use cert and key from config
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}

func localListenAndServe(r *chi.Mux) {
	http.ListenAndServe(":5120", r)
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
