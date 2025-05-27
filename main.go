package main

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
	"github.com/joho/godotenv"
	"github.com/melkdesousa/gamgo/dao"
	"github.com/melkdesousa/gamgo/database"
	_ "github.com/melkdesousa/gamgo/docs/swagger" // This is required for swagger to work
	"github.com/melkdesousa/gamgo/external/rawg"
	"github.com/melkdesousa/gamgo/handlers"
	"github.com/melkdesousa/gamgo/services"
)

func init() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found or error loading .env file")
	}
}

//	@title						Gamgo API
//	@version					1.0
//	@description				Gamgo is a game search API that allows users to search for games by title, leveraging both local database and external APIs.
//	@license.name				Apache 2.0
//	@license.url				http://www.apache.org/licenses/LICENSE-2.0.html
//	@host						localhost:3000
//	@BasePath					/
//	@schemes					http
//	@externalDocs.description	OpenAPI
func main() {
	// Initialize environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found or error loading .env file")
	}
	// Initialize Database Connection
	dbConn := database.GetDBConnection()
	defer func() {
		if err := dbConn.Close(context.Background()); err != nil {
			log.Fatalf("Failed to close database connection: %v", err)
		}
		log.Println("Database connection closed.")
	}()

	// Initialize Cache Connection
	cacheClient := database.GetCacheConnection()
	defer func() {
		if err := cacheClient.Close(); err != nil {
			log.Fatalf("Failed to close cache connection: %v", err)
		}
		log.Println("Cache connection closed.")
	}()

	// Initialize DAO
	gameDAO := dao.NewGameDAO(dbConn)

	// Initialize External APIs
	rawgAPI := rawg.NewRawgAPI()

	// Initialize Services
	gameService := services.NewGameService(gameDAO, cacheClient, rawgAPI)

	// Initialize Handlers
	gameHandler := handlers.NewGameHandler(gameService)

	// Initialize Fiber App
	app := fiber.New(fiber.Config{
		AppName: "gamgo",
	})

	app.Get("/swagger/*", swagger.HandlerDefault) // default

	// Middleware
	app.Use(recover.New())

	// Routes
	app.Get("/games/search", gameHandler.SearchGames)
	app.Get("/games", gameHandler.ListGames)
	// Potentially other routes for health checks, etc.
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	// Start Server
	log.Println("Starting server on port :3000")
	if err := app.Listen(":3000"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
