package database_test

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"

	"github.com/finlleyl/shorty_reborn/internal/database"
)

func TestExists(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()
	repo := database.NewURLRepository(sqlx.NewDb(db, "sqlmock"))
	ctx := context.Background()

	t.Run("exists true", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT EXISTS (")).
			WithArgs("alias123").
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

		ok, err := repo.Exists(ctx, "alias123")
		require.NoError(t, err)
		require.True(t, ok)
	})

	t.Run("exists false", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT EXISTS (")).
			WithArgs("alias123").
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

		ok, err := repo.Exists(ctx, "alias123")
		require.NoError(t, err)
		require.False(t, ok)
	})

	t.Run("db error", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT EXISTS (")).
			WithArgs("alias123").
			WillReturnError(errors.New("db error"))

		ok, err := repo.Exists(ctx, "alias123")
		require.Error(t, err)
		require.False(t, ok)
	})

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestSave(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()
	repo := database.NewURLRepository(sqlx.NewDb(db, "sqlmock"))
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO url (alias, url)
		VALUES ($1, $2)
		RETURNING id;`)).
			WithArgs("alias", "http://example.com").
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(10))

		entity, err := repo.Save(ctx, "alias", "http://example.com")
		require.NoError(t, err)
		require.Equal(t, int64(10), entity.ID)
		require.Equal(t, "alias", entity.Alias)
		require.Equal(t, "http://example.com", entity.URL)
	})

	t.Run("scan error", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO url (alias, url)
		VALUES ($1, $2)
		RETURNING id;`)).
			WithArgs("alias", "http://example.com").
			WillReturnRows(sqlmock.NewRows([]string{"id"}))
		_, err := repo.Save(ctx, "alias", "http://example.com")
		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to save url")
	})

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGet(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()
	repo := database.NewURLRepository(sqlx.NewDb(db, "sqlmock"))
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "alias", "url"}).
			AddRow(5, "alias", "http://example.com")
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, alias, url
		FROM url
		WHERE alias = $1;`)).
			WithArgs("alias").
			WillReturnRows(rows)

		entity, err := repo.Get(ctx, "alias")
		require.NoError(t, err)
		require.Equal(t, int64(5), entity.ID)
		require.Equal(t, "alias", entity.Alias)
		require.Equal(t, "http://example.com", entity.URL)
	})

	t.Run("not found", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, alias, url
		FROM url
		WHERE alias = $1;`)).
			WithArgs("alias").
			WillReturnError(sql.ErrNoRows)

		_, err := repo.Get(ctx, "alias")
		require.ErrorIs(t, err, database.ErrNotFound)
	})

	t.Run("db error", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, alias, url
		FROM url
		WHERE alias = $1;`)).
			WithArgs("alias").
			WillReturnError(errors.New("oh no"))

		_, err := repo.Get(ctx, "alias")
		require.Error(t, err)
		require.Contains(t, err.Error(), "postgresURLRepository.Get")
	})

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestDelete(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()
	repo := database.NewURLRepository(sqlx.NewDb(db, "sqlmock"))
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM url
		WHERE alias = $1;`)).
			WithArgs("alias").
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.Delete(ctx, "alias")
		require.NoError(t, err)
	})

	t.Run("not found", func(t *testing.T) {
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM url
		WHERE alias = $1;`)).
			WithArgs("alias").
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := repo.Delete(ctx, "alias")
		require.ErrorIs(t, err, database.ErrNotFound)
	})

	t.Run("exec error", func(t *testing.T) {
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM url
		WHERE alias = $1;`)).
			WithArgs("alias").
			WillReturnError(errors.New("exec fail"))

		err := repo.Delete(ctx, "alias")
		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to delete url")
	})

	t.Run("rows affected error", func(t *testing.T) {
		result := sqlmock.NewErrorResult(errors.New("nope"))
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM url
		WHERE alias = $1;`)).
			WithArgs("alias").
			WillReturnResult(result)

		err := repo.Delete(ctx, "alias")
		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to get rows affected")
	})

	require.NoError(t, mock.ExpectationsWereMet())
}
