package main

import (
	"database/sql"
	"log"
	"log/slog"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"

	"github.com/charlires/go-base-backend-service/internal/adapters/input/rest"
	"github.com/charlires/go-base-backend-service/internal/adapters/input/rest/gen"
	"github.com/charlires/go-base-backend-service/internal/adapters/output"
	"github.com/charlires/go-base-backend-service/internal/core/services"
)

// Log configuration constants. Adjust these to change the logging behaviour
// without touching business logic.
const (
	logLevel  = slog.LevelDebug
	logFormat = "json" // "json" | "text"
)

// newLogger builds a *slog.Logger from the package-level log config constants.
func newLogger() *slog.Logger {
	opts := &slog.HandlerOptions{Level: logLevel}
	if logFormat == "text" {
		return slog.New(slog.NewTextHandler(os.Stdout, opts))
	}
	return slog.New(slog.NewJSONHandler(os.Stdout, opts))
}

func main() {
	// Configure structured logging from constants defined above.
	appLogger := newLogger()
	slog.SetDefault(appLogger)

	// third party dependencies (e.g., database connections, external APIs) would be initialized here

	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	userRepository := output.NewUserSQLiteAdapter(db)
	trackRepository := output.NewTrackSQLiteAdapter(db)
	playlistRepository := output.NewPlaylistSQLiteAdapter(db)

	userService := services.NewUserService(userRepository)
	trackService := services.NewTrackService(trackRepository)
	playlistService := services.NewPlaylistService(playlistRepository, userRepository)

	handlers := rest.NewHandlers(userService, trackService, playlistService)

	slog.Info("starting server on :8080")
	if err := http.ListenAndServe(":8080", gen.Handler(gen.NewStrictHandler(handlers, []gen.StrictMiddlewareFunc{
		rest.NewLoggingMiddleware(appLogger),
	}))); err != nil {
		log.Fatal(err)
	}
}
