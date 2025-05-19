package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/melkdesousa/gamgo/dao"
	"github.com/melkdesousa/gamgo/dao/models"
	"github.com/melkdesousa/gamgo/database"
	"github.com/melkdesousa/gamgo/external/rawg"
	"github.com/redis/go-redis/v9"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func init() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found or error loading .env file")
	}
}

func main() {
	db := database.GetDBConnection()
	cache := database.GetCacheConnection()
	app := fiber.New(fiber.Config{
		AppName: "gamgo",
	})
	app.Use(recover.New())
	defer func() {
		if err := db.Close(context.Background()); err != nil {
			log.Fatalf("Failed to close database connection: %v", err)
		}
	}()
	defer func() {
		if err := cache.Close(); err != nil {
			log.Fatalf("Failed to close cache connection: %v", err)
		}
	}()
	gameDao := dao.NewGameDAO(db)
	rawgAPI := rawg.NewRawgAPI()
	app.Get("/games/search", func(c *fiber.Ctx) error {
		ctx := c.Context()
		title := sanitize(c.Query("title"))
		pageStr := c.Query("page")
		page, err := strconv.Atoi(pageStr)
		if err != nil {
			page = 1
		}
		if page < 1 {
			page = 1
		}
		if title == "" {
			return c.Status(400).JSON(fiber.Map{
				"error": "Title query parameter is required",
			})
		}
		gamesCached, err := cache.Get(
			ctx,
			database.GetCacheKey(
				database.CACHE_SEARCH_GAME_KEY_PREFIX,
				title,
				pageStr,
			),
		).Result()
		if err != nil {
			if err != redis.Nil {
				return c.Status(500).JSON(fiber.Map{
					"error": "Failed to fetch from cache: " + err.Error(),
				})
			}
		}
		if gamesCached != "" {
			var games []models.Game
			if err := json.Unmarshal([]byte(gamesCached), &games); err != nil {
				return c.Status(500).JSON(fiber.Map{
					"error": "Failed to unmarshal cached games: " + err.Error(),
				})
			}
			log.Println("Cache hit, returning cached games")
			return c.JSON(mapGamesModelToJSON(games))
		}
		gamesInDB, err := gameDao.SearchGames(ctx, title)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "Failed to search games in database: " + err.Error(),
			})
		}
		if len(gamesInDB) > 0 {
			gamesJSON, err := SerializerJSON(gamesInDB)
			if err != nil {
				return c.Status(500).JSON(fiber.Map{
					"error": "Failed to serialize game: " + err.Error(),
				})
			}
			if err := cache.Set(ctx,
				database.GetCacheKey(
					database.CACHE_SEARCH_GAME_KEY_PREFIX,
					title,
					pageStr,
				),
				gamesJSON.String(),
				24*time.Hour,
			).Err(); err != nil {
				return c.Status(500).JSON(fiber.Map{
					"error": "Failed to set cache: " + err.Error(),
				})
			}
			return c.JSON(mapGamesModelToJSON(gamesInDB))
		}
		resp, err := rawgAPI.SearchGames(ctx, title, page)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "Failed to search games in external API: " + err.Error(),
			})
		}
		if resp == nil || resp.Count == 0 {
			return c.Status(404).JSON(fiber.Map{
				"error": "No games found",
			})
		}
		games := make([]rawg.Result, len(resp.Results))
		for i, result := range resp.Results {
			games[i] = result
		}
		gamesModel := mapGamesJSONToModel(games)
		gamesJSON, err := SerializerJSON(gamesModel)
		if err != nil {
			log.Println("Error serializing games:", err)
			return c.Status(400).JSON(fiber.Map{
				"error": "Failed to serialize games: " + err.Error(),
			})
		}
		if err := cache.Set(
			ctx,
			database.GetCacheKey(database.CACHE_SEARCH_GAME_KEY_PREFIX, title, pageStr),
			gamesJSON.String(),
			24*time.Hour,
		).Err(); err != nil {
			log.Println("Error setting cache:", err)
			return c.Status(400).JSON(fiber.Map{
				"error": "Failed to set cache: " + err.Error(),
			})
		}
		if err := gameDao.InsertManyGames(ctx, mapGamesJSONToModel(games)); err != nil {
			log.Println("Error inserting games into database:", err)
		}
		return c.Status(200).JSON(mapGamesModelToJSON(gamesModel))
	})
	app.Listen(":3000")
}

func DeserializerJSON[T any](body io.Reader) (T, error) {
	var result T
	err := json.NewDecoder(body).Decode(&result)
	return result, err
}

func SerializerJSON[T any](data T) (*bytes.Buffer, error) {
	var buffer bytes.Buffer
	err := json.NewEncoder(&buffer).Encode(data)
	return &buffer, err
}

func mapGameModelToJSON(game models.Game) map[string]interface{} {
	gameMap := map[string]interface{}{
		"id":        game.ID,
		"title":     game.Title,
		"released":  game.ReleaseDate,
		"platforms": game.Platforms,
		"rating":    game.Rating,
		// "description": game.Description,
	}
	return gameMap
}

func mapGamesModelToJSON(games []models.Game) []map[string]interface{} {
	gamesMap := make([]map[string]interface{}, len(games))
	for i, game := range games {
		gamesMap[i] = mapGameModelToJSON(game)
	}
	return gamesMap
}

func mapGameJSONExternalToModel(gameJSON rawg.Result) models.Game {
	platforms := make([]string, len(gameJSON.Platforms))
	for i, p := range gameJSON.Platforms {
		platforms[i] = p.Platform.Name
	}
	releaseDate, err := time.Parse(time.DateOnly, gameJSON.Released)
	if err != nil {
		log.Println("Error parsing release date:", err)
		releaseDate = time.Time{}
	}
	game := models.Game{
		ID:    uuid.NewString(),
		Title: gameJSON.Name,
		// Description: gameJSON.Description,
		ReleaseDate:    releaseDate.Format(time.DateOnly),
		Platforms:      platforms,
		Rating:         int(gameJSON.Rating * 100),
		ExternalID:     strconv.Itoa(gameJSON.ID),
		CoverImage:     gameJSON.BackgroundImage,
		ExternalSource: gameJSON.Slug,
	}
	return game
}

func mapGamesJSONToModel(games []rawg.Result) []models.Game {
	gamesModel := make([]models.Game, len(games))
	for i, game := range games {
		gamesModel[i] = mapGameJSONExternalToModel(game)
	}
	return gamesModel
}

func sanitize(input string) string {
	input = strings.TrimSpace(input)
	re := regexp.MustCompile(`[^a-zA-Z0-9\s]`)
	return re.ReplaceAllString(input, "")
}
