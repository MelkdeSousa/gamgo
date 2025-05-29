package services

import (
	"context"
	"time"

	"github.com/melkdesousa/gamgo/dao/models"
	"github.com/melkdesousa/gamgo/external/rawg"
	"github.com/redis/go-redis/v9"
)

type GameDAO interface {
	SearchGames(ctx context.Context, title string) ([]models.Game, error)
	InsertManyGames(ctx context.Context, games []models.Game) error
	ListGames(ctx context.Context, page int, platforms []string, title string) ([]models.Game, int, error)
}

type RawgAPI interface {
	SearchGames(ctx context.Context, query string, page int) (*rawg.GameListResponse, error)
}

type Cache interface {
	Get(ctx context.Context, key string) *redis.StringCmd
	Set(ctx context.Context, key string, value any, expiration time.Duration) *redis.StatusCmd
}
