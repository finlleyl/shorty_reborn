package httpserver

import (
	"net/http"

	"github.com/finlleyl/shorty_reborn/internal/config"
)

func NewServer(cfg *config.HTTPServer, handler http.Handler) *http.Server { 
	return &http.Server{
		Addr: cfg.Address,
		Handler: handler,
		ReadTimeout: cfg.Timeout,
		WriteTimeout: cfg.Timeout,
		IdleTimeout: cfg.IdleTimeout,
	}
}