package mappers

import (
	"log"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/melkdesousa/gamgo/dao/models"
	"github.com/melkdesousa/gamgo/external/rawg"
)

// MapGameModelToJSON converts a single Game model to a map for JSON response.
func MapGameModelToJSON(game models.Game) map[string]interface{} {
	gameMap := map[string]interface{}{
		"id":        game.ID,
		"title":     game.Title,
		"released":  game.ReleaseDate,
		"platforms": game.Platforms,
		"rating":    game.Rating,
	}
	return gameMap
}

// MapGamesModelToJSON converts a slice of Game models to a slice of maps for JSON response.
func MapGamesModelToJSON(games []models.Game) []map[string]interface{} {
	gamesMap := make([]map[string]interface{}, len(games))
	for i, game := range games {
		gamesMap[i] = MapGameModelToJSON(game)
	}
	return gamesMap
}

// MapGameJSONExternalToModel converts a game result from the external RAWG API to our internal Game model.
func MapGameJSONExternalToModel(gameJSON rawg.Result) models.Game {
	platforms := make([]string, len(gameJSON.Platforms))
	for i, p := range gameJSON.Platforms {
		platforms[i] = p.Platform.Name
	}

	var releaseDateParsed time.Time
	if gameJSON.Released != "" {
		parsedTime, err := time.Parse(time.DateOnly, gameJSON.Released)
		if err != nil {
			log.Printf("Warning: Error parsing release date '%s': %v. Falling back to empty date.", gameJSON.Released, err)
			releaseDateParsed = time.Time{} // Zero value for time
		} else {
			releaseDateParsed = parsedTime
		}
	}

	return models.Game{
		ID:             uuid.NewString(), // Generate a new UUID for our internal storage
		Title:          gameJSON.Name,
		ReleaseDate:    releaseDateParsed,
		Platforms:      platforms,
		Rating:         int(gameJSON.Rating * 100), // Example: convert 4.5 to 450
		ExternalID:     strconv.Itoa(gameJSON.ID),
		CoverImage:     gameJSON.BackgroundImage,
		ExternalSource: gameJSON.Slug,
	}
}

// MapGamesJSONToModel converts a slice of game results from the external RAWG API to a slice of our internal Game models.
func MapGamesJSONToModel(games []rawg.Result) []models.Game {
	gamesModel := make([]models.Game, len(games))
	for i, game := range games {
		gamesModel[i] = MapGameJSONExternalToModel(game)
	}
	return gamesModel
}
