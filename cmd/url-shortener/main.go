package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/finlleyl/shorty_reborn/internal/config"
	"github.com/finlleyl/shorty_reborn/internal/database"
	"github.com/finlleyl/shorty_reborn/internal/handlers"
	"github.com/finlleyl/shorty_reborn/internal/httpserver"
	"github.com/finlleyl/shorty_reborn/internal/logger"
	"github.com/finlleyl/shorty_reborn/internal/service"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

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

	srv := httpserver.NewServer(&cfg.HTTPServer, r)


	
	go func() {
		ctx, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
		defer cancel()

		<-ctx.Done()
		logger.Info("Shutting down server...")
	}()

	g, gCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		<-gCtx.Done()
		ctxTimeout, cancelTimeout := context.WithTimeout(context.Background(), 15 * time.Second)
		defer cancelTimeout()

		return srv.Shutdown(ctxTimeout)
	})

	g.Go(func() error {
		logger.Infof("Starting server on %s", cfg.HTTPServer.Address)
		return srv.ListenAndServe()
	})

	if err := g.Wait(); err != nil {
		logger.Fatal("Server stopped: %s", err)
	}
}
