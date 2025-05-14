package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"net/url"
	"regexp"

	"github.com/finlleyl/shorty_reborn/internal/database"
)

var (
	ErrInvalidURL   = errors.New("invalid URL")
	ErrInvalidAlias = errors.New("invalid alias")
	ErrAliasExists  = errors.New("alias already exists")
	ErrURLNotFound  = database.ErrNotFound
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

func (s *urlService) Create(ctx context.Context, RawURL, alias string) (*URL, error) {
	parsed, err := url.ParseRequestURI(RawURL)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrInvalidURL, err)
	}

	if alias == "" {
		alias = generateAlias()
	} else { 
		if !isValidAlias(alias) {
			return nil, ErrInvalidAlias
		}
	}

	exists, err := s.repo.Exists(ctx, alias)
	if err != nil {
		return nil, fmt.Errorf("failed to check if alias exists: %s", err)
	}
	if exists {
		return nil, ErrAliasExists
	}

	u, err := s.repo.Save(ctx, alias, parsed.String())
	if err != nil {
		return nil, fmt.Errorf("failed to save url: %s", err)
	}

	return &URL{
		Alias:   u.Alias,
		OrigURL: u.URL,
	}, nil
}

func (s *urlService) Resolve(ctx context.Context, alias string) (*URL, error) {
    u, err := s.repo.Get(ctx, alias)
    if err != nil {
        if errors.Is(err, database.ErrNotFound) {
            return nil, fmt.Errorf("resolve: %w", ErrURLNotFound)
        }
        return nil, fmt.Errorf("resolve: %w", err)
    }

    return &URL{
        Alias:   u.Alias,
        OrigURL: u.URL,
    }, nil
}

func (s *urlService) Delete(ctx context.Context, alias string) error {
	err := s.repo.Delete(ctx, alias)
	if err != nil {
		switch {
			case errors.Is(err, database.ErrNotFound):
				return fmt.Errorf("delete: %w", ErrURLNotFound)
			default:
				return fmt.Errorf("delete: %w", err)
		}
	}
	
	return nil
}

func generateAlias() string {
	b := make([]byte, 6)
	_, _ = rand.Read(b)
	s := base64.RawURLEncoding.EncodeToString(b)
	if len(s) > 6 {
		return s[:6]
	}

	return s
}

var aliasRegexp = regexp.MustCompile(`^[A-Za-z0-9_-]{3,10}$`)

func isValidAlias(alias string) bool {
	return aliasRegexp.MatchString(alias)
}
