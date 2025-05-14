package main 

import (
	"fmt"
	"log"

	"github.com/finlleyl/shorty_reborn/internal/config"
	"github.com/finlleyl/shorty_reborn/internal/database"
	"github.com/finlleyl/shorty_reborn/internal/logger"
)

func main() {
	cfg := config.MustLoad()
	fmt.Println(cfg)

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

	logger.Info("Starting URL shortener")
}