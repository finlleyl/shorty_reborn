package main 

import (
	"context"
	"log"

	"github.com/finlleyl/shorty_reborn/internal/config"
	"github.com/finlleyl/shorty_reborn/internal/database"
	"github.com/finlleyl/shorty_reborn/internal/logger"
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

	ctx := context.Background()
	urlRepo := database.NewURLRepository(db)
	
	url, err := urlRepo.Save(ctx, "test", "https://google.com")
	if err != nil {
		logger.Fatal("Failed to save url: %s", err)
	}
	logger.Info("URL saved: %s", url)

	logger.Info("Starting URL shortener")
}