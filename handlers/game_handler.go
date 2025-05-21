package handlers

import (
    "log"
    "net/http"
    "strconv"
    "strings"

    "github.com/gofiber/fiber/v2"
    "github.com/melkdesousa/gamgo/mappers"
    "github.com/melkdesousa/gamgo/services"
    "github.com/melkdesousa/gamgo/utils"
)

// GameHandler handles HTTP requests related to games.
type GameHandler struct {
    gameService *services.GameService
}

// NewGameHandler creates a new GameHandler.
func NewGameHandler(gameService *services.GameService) *GameHandler {
    return &GameHandler{
        gameService: gameService,
    }
}

// SearchGames handles the /games/search endpoint.
func (h *GameHandler) SearchGames(c *fiber.Ctx) error {
    ctx := c.Context()
    titleQuery := c.Query("title")
    pageStr := c.Query("page", "1") // Default page to "1"

    page, err := strconv.Atoi(pageStr)
    if err != nil || page < 1 {
        page = 1 // Default to page 1 if conversion fails or page is invalid
    }

    sanitizedTitle := utils.Sanitize(titleQuery)

    if sanitizedTitle == "" {
        return c.Status(http.StatusBadRequest).JSON(fiber.Map{
            "error": "Title query parameter is required and cannot be empty after sanitization",
        })
    }

    log.Printf("Handler: Searching games with title (sanitized): '%s', page: %d", sanitizedTitle, page)

    games, err := h.gameService.SearchGames(ctx, sanitizedTitle, page, pageStr)
    if err != nil {
        // Determine appropriate HTTP status code based on error type
        // This is a simplified error handling. Production code might inspect errors more deeply.
        log.Printf("Error from GameService: %v", err)
        if strings.Contains(err.Error(), "failed to fetch from cache") ||
            strings.Contains(err.Error(), "failed to unmarshal cached games") ||
            strings.Contains(err.Error(), "failed to search games in database") ||
            strings.Contains(err.Error(), "failed to search games in external API") {
            return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
                "error": err.Error(),
            })
        }
        // Generic internal server error for other unhandled cases
        return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
            "error": "An unexpected error occurred",
        })
    }

    if len(games) == 0 {
        return c.Status(http.StatusNotFound).JSON(fiber.Map{
            "message": "No games found matching your criteria",
        })
    }

    return c.Status(http.StatusOK).JSON(mappers.MapGamesModelToJSON(games))
}