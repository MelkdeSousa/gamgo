package dao

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/melkdesousa/gamgo/dao/models"
)

type GameDAO struct {
	connection *pgx.Conn
}

func NewGameDAO(connection *pgx.Conn) *GameDAO {
	return &GameDAO{
		connection: connection,
	}
}

func (dao *GameDAO) SearchGames(ctx context.Context, term string) ([]models.Game, error) {
	var games []models.Game
	query := `
		SELECT id, title, description, platforms, releaseDate, rating, coverImage, externalId, externalSource FROM games
		WHERE search_vector @@ to_tsquery('english', $1)`
	rows, err := dao.connection.Query(ctx, query, term)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var game models.Game
		if err := rows.Scan(
			&game.ID,
			&game.Title,
			&game.Description,
			&game.ReleaseDate,
			&game.Platforms,
			&game.Rating,
			&game.CoverImage,
			&game.ExternalID,
			&game.ExternalSource,
		); err != nil {
			return nil, err
		}
		games = append(games, game)
	}

	return games, nil
}

func (dao *GameDAO) HasGame(ctx context.Context, term string) (bool, error) {
	var count int
	query := `
		SELECT COUNT(*) FROM games
		WHERE search_vector @@ to_tsquery('english', $1)`
	err := dao.connection.QueryRow(ctx, query, term).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (dao *GameDAO) InsertManyGames(ctx context.Context, games []models.Game) error {
	tx, err := dao.connection.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	for _, game := range games {
		query := `
			INSERT INTO games (title, description, platforms, releaseDate, rating, coverImage, externalId, externalSource)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			ON CONFLICT (externalId) DO NOTHING`
		_, err := tx.Exec(ctx, query, game.Title, game.Description, game.Platforms, game.ReleaseDate, game.Rating, game.CoverImage, game.ExternalID, game.ExternalSource)
		if err != nil {
			return err
		}
	}
	if err := tx.Commit(ctx); err != nil {
		return err
	}
	return nil
}
