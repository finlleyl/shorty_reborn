package service_test

import (
	"context"
	"github.com/finlleyl/shorty_reborn/internal/database"
	"github.com/finlleyl/shorty_reborn/internal/service/servicetest"
	"go.uber.org/mock/gomock"

	"testing"

	"github.com/finlleyl/shorty_reborn/internal/service"
	"github.com/stretchr/testify/require"
)

func TestCreate(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	ctx := context.Background()
	repo := servicetest.NewMockURLRepository(ctrl)
	s := service.NewURLService(repo)

	t.Run("valid url & generated alias", func(t *testing.T) {
		url := "https://golang.org"
		repo.EXPECT().
			Exists(ctx, gomock.Any()).
			Return(false, nil)
		repo.EXPECT().
			Save(ctx, gomock.Any(), url).
			Return(&database.URL{ID: 1, Alias: "abc123", URL: url}, nil)

		got, err := s.Create(ctx, url, "")
		require.NoError(t, err)
		require.Equal(t, "abc123", got.Alias)
		require.Equal(t, url, got.OrigURL)
	})
}
