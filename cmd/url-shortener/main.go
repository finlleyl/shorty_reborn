package main 

import (
	"log"
	"net/http"

	"github.com/finlleyl/shorty_reborn/internal/config"
	"github.com/finlleyl/shorty_reborn/internal/database"
	"github.com/finlleyl/shorty_reborn/internal/handlers"
	"github.com/finlleyl/shorty_reborn/internal/httpserver"
	"github.com/finlleyl/shorty_reborn/internal/logger"
	"github.com/finlleyl/shorty_reborn/internal/service"
)

func main() {
	cfg := config.MustLoad()

	logger, cleanup, err := logger.NewSugared(logger.Mode(cfg.Env))
	if err != nil {
		log.Fatalf("Failed to create logger: %s", err)
	}
	logger.Info("Logger created")
	defer cleanup()

	db, err := database.NewDB(&cfg.Database)
	if err != nil {
		logger.Fatal("Failed to create db: %s", err)
	}
	logger.Info("DB created")
	defer db.Close()

	urlRepo := database.NewURLRepository(db)
	
	urlService := service.NewURLService(urlRepo)
	handler := handlers.NewHandler(urlService)

	r := httpserver.NewRouter(handler, logger)

	srv := &http.Server{
		Addr: cfg.HTTPServer.Address,
		Handler: r,
		ReadTimeout: cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout: cfg.HTTPServer.IdleTimeout,
	}

	logger.Infof("Starting server on %s", cfg.HTTPServer.Address)
	if err := srv.ListenAndServe(); err != nil {
		logger.Fatal("Failed to start server: %s", err)
	}
}