package main

import (
	"database/sql"
	"log"
	"log/slog"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"

	"github.com/charlires/go-base-backend-service/internal/adapters/input"
	"github.com/charlires/go-base-backend-service/internal/adapters/output"
	"github.com/charlires/go-base-backend-service/internal/core/services"
)

func main() {
	// Configure structured logging — JSON format, debug level, output to stdout
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})))

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

	userHandler := input.NewUserHTTPAdapter(userService)
	trackHandler := input.NewTrackHTTPAdapter(trackService)
	playlistHandler := input.NewPlaylistHTTPAdapter(playlistService)

	http.HandleFunc("/users/{userid}", userHandler.GetUserByID)
	http.HandleFunc("/tracks/{trackid}", trackHandler.GetTrackByID)
	http.HandleFunc("/playlists/{playlistid}", playlistHandler.GetPlaylistByID)
	http.HandleFunc("/playlists", playlistHandler.CreatePlaylist)
	http.HandleFunc("/playlists/{playlistid}/tracks", playlistHandler.AddTrackToPlaylist)

	// go func() {
	slog.Info("starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
	// }()
	// println("Hello, World!")
}
