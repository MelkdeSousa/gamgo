package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	_ "github.com/melkdesousa/gamgo/docs/swagger" // This is required for swagger to work
)

type SwaggerHandler struct {
	app *fiber.App
}

func NewSwaggerHandler(app *fiber.App) {
	handler := &SwaggerHandler{
		app: app,
	}
	handler.setupRoutes()
	handler.app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})
}

// @title						Gamgo API
// @version					1.0
// @description				Gamgo is a game search API that allows users to search for games by title, leveraging both local database and external APIs.
// @license.name				Apache 2.0
// @license.url				http://www.apache.org/licenses/LICENSE-2.0.html
// @host						localhost:3000
// @BasePath					/
// @schemes					http
// @externalDocs.description	OpenAPI
//
// @securityDefinitions.apikey	JWT
// @in							header
// @name						Authorization
// @description				Enter your JWT token in the format: Bearer \<token\>
func (h *SwaggerHandler) setupRoutes() {
	h.app.Get("/swagger/*", swagger.HandlerDefault) // default
}
