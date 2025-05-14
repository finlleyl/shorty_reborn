package main 

import (
	"fmt"
	"log"

	"github.com/finlleyl/shorty_reborn/internal/config"
	"github.com/finlleyl/shorty_reborn/internal/logger"
)

func main() {
	cfg := config.MustLoad()
	fmt.Println(cfg)

	logger, cleanup, err := logger.NewSugared(logger.Mode(cfg.Env))
	if err != nil {
		log.Fatalf("Failed to create logger: %s", err)
	}
	defer cleanup()

	logger.Info("Starting URL shortener")
}