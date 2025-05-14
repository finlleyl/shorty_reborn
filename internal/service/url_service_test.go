package service_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/finlleyl/shorty_reborn/internal/database"
	"github.com/finlleyl/shorty_reborn/internal/service"
	"github.com/finlleyl/shorty_reborn/internal/service/servicetest"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestCreate_AllCases(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	repo := servicetest.NewMockURLRepository(ctrl)
	svc := service.NewURLService(repo)

	t.Run("invalid URL", func(t *testing.T) {
		_, err := svc.Create(ctx, "%%%://bad-url", "")
		require.ErrorIs(t, err, service.ErrInvalidURL)
	})

	t.Run("invalid alias", func(t *testing.T) {
		_, err := svc.Create(ctx, "https://valid.com", "no spaces")
		require.ErrorIs(t, err, service.ErrInvalidAlias)
	})

	t.Run("exists check error", func(t *testing.T) {
		raw := "https://ok.com"
		repo.EXPECT().
			Exists(ctx, gomock.Any()).
			Return(false, fmt.Errorf("db down"))
		_, err := svc.Create(ctx, raw, "")
		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to check if alias exists")
	})

	t.Run("alias already exists", func(t *testing.T) {
		raw := "https://ok.com"
		repo.EXPECT().
			Exists(ctx, gomock.Any()).
			Return(true, nil)
		_, err := svc.Create(ctx, raw, "foo123")
		require.ErrorIs(t, err, service.ErrAliasExists)
	})

	t.Run("save error", func(t *testing.T) {
		raw := "https://ok.com"
		validAlias := "alias1"
		repo.EXPECT().
			Exists(ctx, validAlias).
			Return(false, nil)
		repo.EXPECT().
			Save(ctx, validAlias, raw).
			Return(nil, fmt.Errorf("write fail"))
		_, err := svc.Create(ctx, raw, validAlias)
		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to save url")
	})

	t.Run("success with provided alias", func(t *testing.T) {
		raw := "https://ok.com"
		given := "myalias"
		repo.EXPECT().
			Exists(ctx, given).
			Return(false, nil)
		repo.EXPECT().
			Save(ctx, given, raw).
			Return(&database.URL{ID: 42, Alias: given, URL: raw}, nil)

		out, err := svc.Create(ctx, raw, given)
		require.NoError(t, err)
		require.Equal(t, given, out.Alias)
		require.Equal(t, raw, out.OrigURL)
	})

	// Пример, как протестировать генерацию случайного alias:
	t.Run("success with generated alias", func(t *testing.T) {
		raw := "https://golang.org"
		// любой alias проходит Exists и Save
		repo.EXPECT().
			Exists(ctx, gomock.Any()).
			Return(false, nil)
		repo.EXPECT().
			Save(ctx, gomock.Any(), raw).
			DoAndReturn(func(_ context.Context, alias, url string) (*database.URL, error) {
				// проверяем, что alias сгенерирован и валиден по regexp
				require.Regexp(t, `^[A-Za-z0-9_-]{6}$`, alias)
				return &database.URL{ID: 1, Alias: alias, URL: url}, nil
			})

		out, err := svc.Create(ctx, raw, "")
		require.NoError(t, err)
		require.Len(t, out.Alias, 6)
		require.Equal(t, raw, out.OrigURL)
	})
}

func TestResolve(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	repo := servicetest.NewMockURLRepository(ctrl)
	svc := service.NewURLService(repo)

	t.Run("not found", func(t *testing.T) {
		repo.EXPECT().
			Get(ctx, "foo").
			Return(nil, database.ErrNotFound)

		_, err := svc.Resolve(ctx, "foo")
		require.ErrorIs(t, err, service.ErrURLNotFound)
	})

	t.Run("db error", func(t *testing.T) {
		repo.EXPECT().
			Get(ctx, "alias").
			Return(nil, fmt.Errorf("oops"))

		_, err := svc.Resolve(ctx, "alias")
		require.Error(t, err)
		require.Contains(t, err.Error(), "resolve:")
	})

	t.Run("success", func(t *testing.T) {
		repo.EXPECT().
			Get(ctx, "good").
			Return(&database.URL{Alias: "good", URL: "https://ok.com"}, nil)

		out, err := svc.Resolve(ctx, "good")
		require.NoError(t, err)
		require.Equal(t, "good", out.Alias)
		require.Equal(t, "https://ok.com", out.OrigURL)
	})
}

func TestDelete(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	repo := servicetest.NewMockURLRepository(ctrl)
	svc := service.NewURLService(repo)

	t.Run("not found", func(t *testing.T) {
		repo.EXPECT().
			Delete(ctx, "missing").
			Return(database.ErrNotFound)

		err := svc.Delete(ctx, "missing")
		require.ErrorIs(t, err, service.ErrURLNotFound)
	})

	t.Run("db error", func(t *testing.T) {
		repo.EXPECT().
			Delete(ctx, "alias").
			Return(fmt.Errorf("cannot delete"))

		err := svc.Delete(ctx, "alias")
		require.Error(t, err)
		require.Contains(t, err.Error(), "delete:")
	})

	t.Run("success", func(t *testing.T) {
		repo.EXPECT().
			Delete(ctx, "foo").
			Return(nil)

		err := svc.Delete(ctx, "foo")
		require.NoError(t, err)
	})
}
