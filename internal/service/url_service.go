package service

import (
	"context"

	"github.com/finlleyl/shorty_reborn/internal/database"
)

type URL struct {
	Alias   string
	OrigURL string
}

type URLService interface {
	Create(ctx context.Context, url, alias string) (*URL, error)
	Resolve(ctx context.Context, alias string) (*URL, error)
	Delete(ctx context.Context, alias string) error
}

type urlService struct {
	repo database.URLRepository
}

func NewURLService(r database.URLRepository) URLService {
	return &urlService{repo: r}
}

func (s *urlService) Create(ctx context.Context, url, alias string) (*URL, error) {
	// TODO: IMPLEMENT
	return nil, nil
}

func (s *urlService) Resolve(ctx context.Context, alias string) (*URL, error) {
	// TODO: IMPLEMENT
	return nil, nil
}

func (s *urlService) Delete(ctx context.Context, alias string) error {
	// TODO: IMPLEMENT
	return nil
}
