package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/melkdesousa/gamgo/dao"
	"github.com/melkdesousa/gamgo/dao/models"
	"github.com/melkdesousa/gamgo/database"
	"github.com/melkdesousa/gamgo/external/rawg"
	"github.com/melkdesousa/gamgo/mappers"
	"github.com/melkdesousa/gamgo/utils"
	"github.com/redis/go-redis/v9"
)

// GameService encapsulates business logic related to games.
type GameService struct {
	gameDAO  *dao.GameDAO
	cache    *redis.Client
	rawgAPI  *rawg.RawgAPI
	cacheTTL time.Duration
}

// NewGameService creates a new GameService.
func NewGameService(gameDAO *dao.GameDAO, cache *redis.Client, rawgAPI *rawg.RawgAPI) *GameService {
	cacheTTLStr := os.Getenv("CACHE_TTL_HOURS")
	cacheTTLHours, err := strconv.Atoi(cacheTTLStr)
	if err != nil || cacheTTLHours <= 0 {
		log.Printf("Warning: CACHE_TTL_HOURS is not set or invalid in GameService. Defaulting to 24 hours. Error: %v", err)
		cacheTTLHours = 24 // Default TTL
	}
	cacheTTLValue := time.Duration(cacheTTLHours) * time.Hour
	return &GameService{
		gameDAO:  gameDAO,
		cache:    cache,
		rawgAPI:  rawgAPI,
		cacheTTL: cacheTTLValue,
	}
}

// SearchGames searches for games based on title and page.
// It checks cache, then database, then external API.
func (s *GameService) SearchGames(ctx context.Context, sanitizedTitle string, page int, pageStr string) ([]models.Game, error) {
	cacheKey := database.GetCacheKey(database.CACHE_SEARCH_GAME_KEY_PREFIX, sanitizedTitle, pageStr)
	// 1. Check Cache
	gamesCached, err := s.cache.Get(ctx, cacheKey).Result()
	if err != nil && err != redis.Nil {
		log.Printf("Error fetching from cache for key %s: %v", cacheKey, err)
		return nil, fmt.Errorf("failed to fetch from cache: %w", err)
	}
	if gamesCached != "" {
		var games []models.Game
		if err := json.Unmarshal([]byte(gamesCached), &games); err != nil {
			log.Printf("Error unmarshalling cached games for key %s: %v", cacheKey, err)
			return nil, fmt.Errorf("failed to unmarshal cached games: %w", err)
		}
		log.Printf("Cache hit for key %s, returning cached games", cacheKey)
		return games, nil
	}
	log.Printf("Cache miss for key %s", cacheKey)
	// 2. Check Database
	// Note: Current DB search in DAO is by a general term, not specifically by page.
	// If pagination from DB is needed, DAO would need adjustment.
	// For now, we assume if any games match the title, we return them and cache them with the pageStr.
	gamesInDB, err := s.gameDAO.SearchGames(ctx, sanitizedTitle)
	if err != nil {
		log.Printf("Error searching games in database for title '%s': %v", sanitizedTitle, err)
		return nil, fmt.Errorf("failed to search games in database: %w", err)
	}
	if len(gamesInDB) > 0 {
		log.Printf("Found %d games in DB for title '%s'", len(gamesInDB), sanitizedTitle)
		gamesJSON, err := utils.SerializerJSON(gamesInDB)
		if err != nil {
			log.Printf("Error serializing games from DB for caching (key %s): %v", cacheKey, err)
			// Non-fatal for returning data, but cache won't be set
		} else {
			if err := s.cache.Set(ctx, cacheKey, gamesJSON.String(), s.cacheTTL).Err(); err != nil {
				log.Printf("Error setting cache for DB results (key %s): %v", cacheKey, err)
				// Non-fatal
			} else {
				log.Printf("Successfully cached DB results for key %s with TTL %v", cacheKey, s.cacheTTL)
			}
		}
		return gamesInDB, nil
	}
	log.Printf("No games found in DB for title '%s'. Fetching from RAWG API.", sanitizedTitle)
	// 3. Fetch from External API (RAWG)
	resp, err := s.rawgAPI.SearchGames(ctx, sanitizedTitle, page)
	if err != nil {
		log.Printf("Error searching games in external API for title '%s', page %d: %v", sanitizedTitle, page, err)
		return nil, fmt.Errorf("failed to search games in external API: %w", err)
	}
	if resp == nil || resp.Count == 0 {
		log.Printf("No games found in external API for title '%s', page %d", sanitizedTitle, page)
		return []models.Game{}, nil // Return empty slice, not an error for "not found"
	}
	gamesFromAPI := make([]rawg.Result, len(resp.Results))
	for i, result := range resp.Results {
		gamesFromAPI[i] = result
	}
	gamesModel := mappers.MapGamesJSONToModel(gamesFromAPI)
	// Cache the results from RAWG API
	gamesJSON, err := utils.SerializerJSON(gamesModel)
	if err != nil {
		log.Printf("Error serializing games from API for caching (key %s): %v", cacheKey, err)
		// Non-fatal for returning data, but cache won't be set
	} else {
		if err := s.cache.Set(ctx, cacheKey, gamesJSON.String(), s.cacheTTL).Err(); err != nil {
			log.Printf("Error setting cache for API results (key %s): %v", cacheKey, err)
			// Non-fatal
		} else {
			log.Printf("Successfully cached API results for key %s with TTL %v", cacheKey, s.cacheTTL)
		}
	}
	// Asynchronously save to DB (or synchronously if critical)
	// For simplicity here, it's synchronous. Consider a background worker for bulk/non-critical inserts.
	if err := s.gameDAO.InsertManyGames(ctx, gamesModel); err != nil {
		log.Printf("Error inserting games from API into database (title '%s'): %v", sanitizedTitle, err)
		// Log error but still return data from API
	} else {
		log.Printf("Successfully inserted %d games from API into database for title '%s'", len(gamesModel), sanitizedTitle)
	}
	return gamesModel, nil
}

// ListGames retrieves a list of games from the database, with optional filters.
func (s *GameService) ListGames(ctx context.Context, page int, platforms []string, title string) (games []models.Game, total int, err error) {
	games, total, err = s.gameDAO.ListGames(ctx, page, platforms, title)
	if err != nil {
		log.Printf("Error listing games from database: %v", err)
		return nil, 0, fmt.Errorf("failed to list games: %w", err)
	}

	if len(games) == 0 {
		log.Println("No games found in database.")
		return nil, 0, nil // No error, just no games found
	}

	return games, total, nil
}
