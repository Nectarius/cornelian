package server

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"

	"time"

	"github.com/golangcollege/sessions"
	gsessions "github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

func SetupSessions() (*sessions.Session, error) {

	// User session mgmt, used once gothic has performed the social login.
	var session *sessions.Session
	var secret = []byte("replace-me-with-something")
	session = sessions.New(secret)
	session.Lifetime = 3 * time.Hour

	return localSetupSessions(session)
}

func SetupSessionsForServerWithHttps(session *sessions.Session) (*sessions.Session, error) {

	// Gothic session mgmt for social login
	key := "Secret-session-key" // Replace with your SESSION_SECRET or similar
	maxAge := 86400 * 30        // 30 days
	isProd := false             // Set to true when serving over https

	cookieStore := gsessions.NewCookieStore([]byte(key))
	cookieStore.MaxAge(maxAge)

	// local
	// cookieStore.Options.Domain = "localhost"
	cookieStore.Options.Domain = "kornelian.com"
	cookieStore.Options.Path = ""
	cookieStore.Options.HttpOnly = true // HttpOnly should always be enabled
	cookieStore.Options.Secure = isProd

	gothic.Store = cookieStore

	// Read Google auth credentials from .credentials file.
	clientID, clientSecret := readCredentials()

	// local : http://localhost:5120/auth/google/callback
	// server  https://kornelian.com/auth/google/callback?provider=google
	// uk server http://80.190.84.21/
	goth.UseProviders(
		google.New(clientID, clientSecret, "https://kornelian.com/auth/google/callback?provider=google", "email", "profile"),
	)

	return session, nil
}

func localSetupSessions(session *sessions.Session) (*sessions.Session, error) {

	// Gothic session mgmt for social login
	key := "Secret-session-key" // Replace with your SESSION_SECRET or similar
	maxAge := 86400 * 30        // 30 days
	isProd := false             // Set to true when serving over https

	cookieStore := gsessions.NewCookieStore([]byte(key))
	cookieStore.MaxAge(maxAge)

	cookieStore.Options.Domain = "localhost"
	cookieStore.Options.Path = ""
	cookieStore.Options.HttpOnly = true // HttpOnly should always be enabled
	cookieStore.Options.Secure = isProd

	gothic.Store = cookieStore

	// Read Google auth credentials from .credentials file.
	clientID, clientSecret := readCredentials()

	// local : http://localhost:5120/auth/google/callback
	// server  https://kornelian.com/auth/google/callback?provider=google
	// uk server http://80.190.84.21/
	goth.UseProviders(
		google.New(clientID, clientSecret, "http://localhost:5120/auth/google/callback", "email", "profile"),
	)

	return session, nil
}

func logoutHandler(session *sessions.Session) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		session.Destroy(r)
		gothic.Logout(w, r)
		w.Header().Set("HX-Redirect", "/") // Use the special HTMX redirect header to trigger a full page reload.
	}
}

func readCredentials() (string, string) {
	//	currentFile, err := os.Executable()
	//    if err != nil {
	//    slog.Error(err.Error())
	//	os.Exit(1)
	//    }
	//    currentDir := filepath.Dir(currentFile)

	// Construct the path to the file you want to read
	//        filePath := filepath.Join(currentDir, "underlying/recources/.credentials")

	credzData, err := os.ReadFile(".credentials")

	//credzData, err := os.ReadFile("underlying/recources/.credentials")
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	credz := make(map[string]interface{})
	if err := json.Unmarshal(credzData, &credz); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	clientID := credz["web"].(map[string]interface{})["client_id"].(string)
	clientSecret := credz["web"].(map[string]interface{})["client_secret"].(string)
	return clientID, clientSecret
}
