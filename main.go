package main

import (
	"context"
	"log"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/html/v2"
	"github.com/joho/godotenv"
	"github.com/melkdesousa/gamgo/config"
	"github.com/melkdesousa/gamgo/dao"
	"github.com/melkdesousa/gamgo/database"
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

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found or error loading .env file")
	}
	dbConn := database.GetDBConnection()
	defer func() {
		if err := dbConn.Close(context.Background()); err != nil {
			log.Fatalf("Failed to close database connection: %v", err)
		}
		log.Println("Database connection closed.")
	}()
	cacheClient := database.GetCacheConnection()
	defer func() {
		if err := cacheClient.Close(); err != nil {
			log.Fatalf("Failed to close cache connection: %v", err)
		}
		log.Println("Cache connection closed.")
	}()
	gameDAO := dao.NewGameDAO(dbConn)
	accountDAO := dao.NewAccountDAO(dbConn)
	rawgAPI := rawg.NewRawgAPI()
	gameService := services.NewGameService(gameDAO, cacheClient, rawgAPI)
	accountService := services.NewAccountService(accountDAO)
	engine := html.New("views", ".html")
	app := fiber.New(fiber.Config{
		AppName: "gamgo",
		Views:   engine,
	})
	app.Static("/static", "./views/static")
	handlers.NewSwaggerHandler(app)
	handlers.NewAuthHandler(app, accountService)
	app.Use(recover.New())
	app.Use(JWTProtection)
	handlers.NewGameHandler(app, gameService)
	log.Println("Starting server on port :3000")
	if err := app.Listen(":3000"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func JWTProtection(c *fiber.Ctx) error {
	return jwtware.New(jwtware.Config{
		SigningKey:  jwtware.SigningKey{Key: []byte(config.MustGetEnv("JWT_SECRET"))},
		TokenLookup: "cookie:token",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			log.Printf("JWT error: %v", err)
			return c.Redirect("/login")
		},
	})(c)
}
