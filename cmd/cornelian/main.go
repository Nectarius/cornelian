package main

import (
	"log/slog"
	"os"

	"github.com/nefarius/cornelian/underlying/app/access"
	"github.com/nefarius/cornelian/underlying/app/server"
	"github.com/nefarius/cornelian/underlying/app/store"
	"golang.org/x/net/context"
)

func main() {

	session, err := server.SetupSessions()
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	var cornelianModule = access.NewCornelianModule()
	var questionRepository = cornelianModule.QuestionRepository
	var quiz = questionRepository.GetQuiz()
	// In-memory DB for questions and answers.
	db := store.NewInMem()
	db.FillWithData(quiz.Questions)

	defer func(module *access.CornelianModule, ctx context.Context) {
		module.Clear()
	}(cornelianModule, context.Background())
	// Setup and start HTTP server
	server.StartServer(session, db, cornelianModule)
}
