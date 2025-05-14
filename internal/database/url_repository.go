package database

import (
	"context"	
	"errors"
	"fmt"
	"database/sql"

	"github.com/jmoiron/sqlx"
)	

type URL struct {
	ID int64 `db:"id"`
	Alias string `db:"alias"`
	URL string `db:"url"`
}

var ErrNotFound = errors.New("url not found")
type URLRepository interface {
	Exists(ctx context.Context, alias string) (bool, error)
	Save(ctx context.Context, alias, url string) (*URL, error)
	Get(ctx context.Context, alias string) (*URL, error)
	Delete(ctx context.Context, alias string) error
}

type postgresURLRepository struct {
	db *sqlx.DB
}

func NewURLRepository(db *sqlx.DB) URLRepository {
	return &postgresURLRepository{db: db}
}

func (r *postgresURLRepository) Exists(ctx context.Context, alias string) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM url
			WHERE alias = $1
		)
	`

	var exists bool
	if err := r.db.GetContext(ctx, &exists, query, alias); err != nil {
		return false, fmt.Errorf("failed to check if alias exists: %w", err)
	}

	return exists, nil
}

func (r *postgresURLRepository) Save(ctx context.Context, alias, url string) (*URL, error) {
	query := `
		INSERT INTO url (alias, url)
		VALUES ($1, $2)
		RETURNING id;
	`

	urlEntity := &URL{
		Alias: alias,
		URL: url,
	}

	row := r.db.QueryRowContext(ctx, query, alias, url) 
	if err := row.Scan(&urlEntity.ID); err != nil {
		return nil, fmt.Errorf("failed to save url: %w", err)
	}
	
	return urlEntity, nil
}

func (r *postgresURLRepository) Get(ctx context.Context, alias string) (*URL, error) {
    query := `
        SELECT id, alias, url
        FROM url
        WHERE alias = $1;
    `
    var urlEntity URL
    err := r.db.GetContext(ctx, &urlEntity, query, alias)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, ErrNotFound
        }
        return nil, fmt.Errorf("postgresURLRepository.Get: %w", err)
    }
    return &urlEntity, nil
}

func (r *postgresURLRepository) Delete(ctx context.Context, alias string) error {
	query := `
		DELETE FROM url
		WHERE alias = $1;
	`

	_, err := r.db.ExecContext(ctx, query, alias)
	if err != nil {
		return fmt.Errorf("failed to delete url: %w", err)
	}

	return nil
}