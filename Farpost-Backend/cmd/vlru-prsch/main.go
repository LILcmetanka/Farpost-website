// @title VLRU-PRSCH API
// @version 1.0
// @description API для системы VLRU-PRSCH

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
package main

import (
	"log/slog"
	"net/http"
	"os"
	"vlru-prsch/internal/config"
	blackoutsget "vlru-prsch/internal/http-server/handlers/blackouts/get"
	orgsget "vlru-prsch/internal/http-server/handlers/organizations/get"
	dayget "vlru-prsch/internal/http-server/handlers/calendar/day/get"
	monthget "vlru-prsch/internal/http-server/handlers/calendar/month/get"
	"vlru-prsch/internal/http-server/handlers/complaints"
	"vlru-prsch/internal/http-server/handlers/search"
	"vlru-prsch/internal/lib/logger/sl"
	"vlru-prsch/internal/lib/logger/slogpretty"
	"vlru-prsch/internal/storage/sqlite"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

    _ "vlru-prsch/docs"
    httpSwagger "github.com/swaggo/http-swagger"
)

const (
	local = "local"
	dev   = "dev"
	prod  = "prod"
)

// main godoc
// @Summary Запуск сервера
// @Description Запускает HTTP сервер с API
func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info("init config and start app", slog.Any("cfg", cfg))

	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
	}

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Use(corsConfig(cfg.Env))

	router.Get("/swagger/*", httpSwagger.Handler(
        httpSwagger.URL("/swagger/doc.json"), 
    ))

	router.Route("/off", func(r chi.Router) {
		r.Post("/search", search.New(log, storage))
		r.Get("/blackouts", blackoutsget.New(log, storage))
		r.Get("/orgs", orgsget.New(log, storage))
		r.Get("/complaints", complaints.New(log, storage))
		r.Get("/calendar", monthget.New(log, storage))
		r.Get("/calendar/day", dayget.New(log, storage))
	})

	log.Info("starting server", slog.Any("address", cfg.Address))

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
		IdleTimeout:  cfg.IddleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}

	log.Error("server stoped")
}

func corsConfig(env string) func(next http.Handler) http.Handler {
	switch env {
	case local:
		return cors.Handler(cors.Options{
			AllowedOrigins:   []string{"http://localhost", "http://127.0.0.1"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Content-Type", "Authorization"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: true,
			MaxAge:           300,
		})
	// dev/prod можно настроить аналогично
	}

	return cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: false,
		MaxAge:           300,
	})
}



func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case local:
		log = setupPrettySlog()
	case dev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case prod:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
