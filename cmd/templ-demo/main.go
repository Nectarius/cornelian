package main

import (
	"log/slog"
	"os"

	"github.com/nefarius/cornelian/internal/app/server"
	"github.com/nefarius/cornelian/internal/app/store"
)

func main() {

	session, err := server.SetupSessions()
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	// In-memory DB for questions and answers.
	db := store.NewInMem()
	db.SeedWithFakeData()

	// Setup and start HTTP server
	server.StartServer(session, db)
}
