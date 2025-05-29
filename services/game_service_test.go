package services

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/melkdesousa/gamgo/dao"
	"github.com/melkdesousa/gamgo/dao/models"
	"github.com/melkdesousa/gamgo/database"
	"github.com/melkdesousa/gamgo/external/rawg"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockGameDAO is a mock implementation of the GameDAO
type MockGameDAO struct {
	mock.Mock
}

func (m *MockGameDAO) SearchGames(ctx context.Context, title string) ([]models.Game, error) {
	args := m.Called(ctx, title)
	return args.Get(0).([]models.Game), args.Error(1)
}

func (m *MockGameDAO) InsertManyGames(ctx context.Context, games []models.Game) error {
	args := m.Called(ctx, games)
	return args.Error(0)
}

func (m *MockGameDAO) ListGames(ctx context.Context, page int, platforms []string, title string) ([]models.Game, int, error) {
	args := m.Called(ctx, page, platforms, title)
	return args.Get(0).([]models.Game), args.Int(1), args.Error(2)
}

// MockRawgAPI is a mock implementation of the RawgAPI
type MockRawgAPI struct {
	mock.Mock
}

func (m *MockRawgAPI) SearchGames(ctx context.Context, query string, page int) (*rawg.GameListResponse, error) {
	args := m.Called(ctx, query, page)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*rawg.GameListResponse), args.Error(1)
}

// MockRedisClient is a mock implementation of the Redis client
type MockRedisClient struct {
	mock.Mock
}

func (m *MockRedisClient) Get(ctx context.Context, key string) *redis.StringCmd {
	args := m.Called(ctx, key)
	return args.Get(0).(*redis.StringCmd)
}

func (m *MockRedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	args := m.Called(ctx, key, value, expiration)
	return args.Get(0).(*redis.StatusCmd)
}

type MockStatusCmd struct {
	mock.Mock
}

func (m *MockStatusCmd) Err() error {
	args := m.Called()
	return args.Error(0)
}

func TestGameService(t *testing.T) {
	err := godotenv.Load("../.env.test")
	assert.NoError(t, err, "Expected no error loading .env file")
	isIntegrationTest := os.Getenv("INTEGRATION_TEST") == "true"
	if !isIntegrationTest {
		t.Skip("Skipping integration test for GameService. Set INTEGRATION_TEST=true to run.")
	}
	gameDAO := dao.NewGameDAO(database.GetDBConnection())
	redisClient := database.GetCacheConnection()
	rawgAPI := rawg.NewRawgAPI()
	gameService := NewGameService(gameDAO, redisClient, rawgAPI)
	ctx := context.Background()
	t.Run("TestSearchGames", func(t *testing.T) {
		games, err := gameService.SearchGames(ctx, "zelda", 1, "1")
		assert.NoError(t, err, "Expected no error when searching for games")
		assert.NotEmpty(t, games, "Expected to find games with title 'zelda'")

		// Verify some properties of the returned games
		for _, game := range games {
			assert.NotEmpty(t, game.Title, "Game title should not be empty")
			assert.NotZero(t, game.ID, "Game ID should not be zero")
		}
	})

	t.Run("TestListGames", func(t *testing.T) {
		games, total, err := gameService.ListGames(ctx, 1, []string{}, "")
		assert.NoError(t, err, "Expected no error when listing games")
		assert.GreaterOrEqual(t, total, 0, "Total should be non-negative")

		if total > 0 {
			assert.NotEmpty(t, games, "Expected to find games when total > 0")

			// Verify some properties of the returned games
			for _, game := range games {
				assert.NotEmpty(t, game.Title, "Game title should not be empty")
				assert.NotZero(t, game.ID, "Game ID should not be zero")
			}
		}
	})

	t.Run("TestListGamesWithPlatformFilter", func(t *testing.T) {
		platforms := []string{"PC", "PlayStation"}
		games, total, err := gameService.ListGames(ctx, 1, platforms, "")
		assert.NoError(t, err, "Expected no error when listing games with platform filter")

		if total > 0 {
			assert.NotEmpty(t, games, "Expected to find games when total > 0")
		}
	})

	t.Run("TestListGamesWithTitleFilter", func(t *testing.T) {
		title := "mario"
		games, total, err := gameService.ListGames(ctx, 1, []string{}, title)
		assert.NoError(t, err, "Expected no error when listing games with title filter")

		if total > 0 {
			assert.NotEmpty(t, games, "Expected to find games when total > 0")
			for _, game := range games {
				assert.Contains(t, game.Title, title, "Game title should contain search term")
			}
		}
	})
}

func TestGameServiceUnit(t *testing.T) {
	// Create mocks
	mockGameDAO := &MockGameDAO{}
	mockRedisClient := &MockRedisClient{}
	mockRawgAPI := &MockRawgAPI{}

	// Create game service with mocks
	gameService := NewGameService(mockGameDAO, mockRedisClient, mockRawgAPI)

	ctx := context.Background()

	t.Run("TestSearchGamesCacheHit", func(t *testing.T) {
		// Setup
		cacheKey := database.GetCacheKey(database.CACHE_SEARCH_GAME_KEY_PREFIX, "mario", "1")
		mockRedisClient.On("Get", ctx, cacheKey).Return(&redis.StringCmd{})

		// Call the service
		games, err := gameService.SearchGames(ctx, "mario", 1, "1")

		// Assertions
		assert.NoError(t, err)
		assert.Len(t, games, 1)
		assert.Equal(t, "Super Mario Bros", games[0].Title)

		// Verify mocks
		mockRedisClient.AssertExpectations(t)
		// The DAO and API should not be called when cache hits
		mockGameDAO.AssertNotCalled(t, "SearchGames")
		mockRawgAPI.AssertNotCalled(t, "SearchGames")
	})

	t.Run("TestSearchGamesDBHit", func(t *testing.T) {
		// Setup
		cacheKey := database.GetCacheKey(database.CACHE_SEARCH_GAME_KEY_PREFIX, "zelda", "1")

		mockRedisClient.On("Get", ctx, cacheKey).Return(&redis.StringCmd{})

		// DB hit
		dbGames := []models.Game{{ID: uuid.NewString(), Title: "Legend of Zelda"}}
		mockGameDAO.On("SearchGames", ctx, "zelda").Return(dbGames, nil)

		// Mock caching DB results
		mockStatusCmd := new(MockStatusCmd)
		mockStatusCmd.On("Err").Return(nil)
		mockRedisClient.On("Set", ctx, cacheKey, mock.Anything, gameService.cacheTTL).Return(mockStatusCmd)

		// Call the service
		games, err := gameService.SearchGames(ctx, "zelda", 1, "1")

		// Assertions
		assert.NoError(t, err)
		assert.Len(t, games, 1)
		assert.Equal(t, "Legend of Zelda", games[0].Title)

		// Verify mocks
		mockRedisClient.AssertExpectations(t)
		mockGameDAO.AssertExpectations(t)
		mockStatusCmd.AssertExpectations(t)
		mockRawgAPI.AssertNotCalled(t, "SearchGames") // API should not be called when DB hits
	})

	t.Run("TestSearchGamesAPIHit", func(t *testing.T) {
		// Setup
		cacheKey := database.GetCacheKey(database.CACHE_SEARCH_GAME_KEY_PREFIX, "metroid", "1")
		mockRedisClient.On("Get", ctx, cacheKey).Return(&redis.StringCmd{})

		// DB miss
		mockGameDAO.On("SearchGames", ctx, "metroid").Return([]models.Game{}, nil)

		// API hit
		apiResponse := &rawg.GameListResponse{
			Count: 1,
			Results: []rawg.Result{
				{
					ID:   3,
					Name: "Metroid Prime",
				},
			},
		}
		mockRawgAPI.On("SearchGames", ctx, "metroid", 1).Return(apiResponse, nil)

		// Mock saving to DB
		mockGameDAO.On("InsertManyGames", ctx, mock.Anything).Return(nil)

		// Mock caching API results
		mockStatusCmd := new(MockStatusCmd)
		mockStatusCmd.On("Err").Return(nil)
		mockRedisClient.On("Set", ctx, cacheKey, mock.Anything, gameService.cacheTTL).Return(mockStatusCmd)

		// Call the service
		games, err := gameService.SearchGames(ctx, "metroid", 1, "1")

		// Assertions
		assert.NoError(t, err)
		assert.Len(t, games, 1)
		assert.Contains(t, games[0].Title, "Metroid")

		// Verify mocks
		mockRedisClient.AssertExpectations(t)
		mockGameDAO.AssertExpectations(t)
		mockRawgAPI.AssertExpectations(t)
		mockStatusCmd.AssertExpectations(t)
	})

	t.Run("TestListGames", func(t *testing.T) {
		// Setup
		page := 1
		platforms := []string{"PC"}
		title := "witcher"

		expectedGames := []models.Game{{ID: uuid.NewString(), Title: "The Witcher 3"}}
		expectedTotal := 1

		mockGameDAO.On("ListGames", ctx, page, platforms, title).Return(expectedGames, expectedTotal, nil)

		// Call the service
		games, total, err := gameService.ListGames(ctx, page, platforms, title)

		// Assertions
		assert.NoError(t, err)
		assert.Equal(t, expectedTotal, total)
		assert.Equal(t, expectedGames, games)

		// Verify mocks
		mockGameDAO.AssertExpectations(t)
	})
}
