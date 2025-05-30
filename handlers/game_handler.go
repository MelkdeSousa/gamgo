package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/melkdesousa/gamgo/mappers"
	"github.com/melkdesousa/gamgo/services"
	"github.com/melkdesousa/gamgo/utils"
)

// GameHandler handles HTTP requests related to games.
type GameHandler struct {
	app         *fiber.App
	gameService *services.GameService
}

// NewGameHandler creates a new GameHandler.
func NewGameHandler(app *fiber.App, gameService *services.GameService) {
	handler := &GameHandler{
		app:         app,
		gameService: gameService,
	}
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("pages/home", fiber.Map{})
	})
	app.Get("/games/search", handler.SearchGames)
	app.Get("/games", handler.ListGames)
}

// SearchGames godoc
//
//	@Summary		Search Games
//	@Description	search games by title
//	@Security		JWT
//	@Tags			games
//	@Accept			json
//	@Produce		json
//	@Param			title	query		string	false	"game search by title"
//	@Param			page	query		int		false	"page number, default is 1"
//	@Success		200		{array}		mappers.PaginationResponse[[]mappers.GameOutputDTO]
//	@Failure		400		{object}	mappers.ErrorResponse
//	@Failure		404		{object}	mappers.PaginationResponse[[]mappers.GameOutputDTO]
//	@Failure		500		{object}	mappers.ErrorResponse
//	@Router			/games/search [get]
func (h *GameHandler) SearchGames(c *fiber.Ctx) error {
	ctx := c.Context()
	titleQuery := c.Query("title", "") // Default to empty string if not provided
	pageStr := c.Query("page", "1")    // Default page to "1"
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1 // Default to page 1 if conversion fails or page is invalid
	}
	sanitizedTitle := utils.Sanitize(titleQuery)
	if sanitizedTitle == "" {
		return c.Status(http.StatusBadRequest).JSON(mappers.ErrorResponse{
			Error:   "Title query parameter is required and cannot be empty after sanitization",
			Details: "Please provide a valid title.",
		})
	}
	log.Printf("Handler: Searching games with title (sanitized): '%s', page: %d", sanitizedTitle, page)
	games, err := h.gameService.SearchGames(ctx, sanitizedTitle, page, pageStr)
	if err != nil {
		log.Printf("Error from GameService: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(mappers.ErrorResponse{
			Error:   "An unexpected error occurred",
			Details: err.Error(),
		})
	}
	if len(games) == 0 {
		return c.Status(http.StatusNotFound).JSON(mappers.PaginationResponse[[]mappers.GameOutputDTO]{
			CommonResponse: mappers.CommonResponse[[]mappers.GameOutputDTO]{
				Data:    []mappers.GameOutputDTO{},
				Message: "No games found matching your criteria",
			},
			Filters: fiber.Map{
				"title": sanitizedTitle,
				"page":  page,
			},
			Page:  page,
			Count: 0,
		})
	}
	return c.Status(http.StatusOK).JSON(mappers.PaginationResponse[[]mappers.GameOutputDTO]{
		CommonResponse: mappers.CommonResponse[[]mappers.GameOutputDTO]{
			Data:    mappers.MapGamesModelToOutputDTO(games),
			Message: "Games retrieved successfully",
		},
		Filters: fiber.Map{
			"title": sanitizedTitle,
			"page":  page,
		},
		Page:  page,
		Count: len(games),
	})
}

// ListGames godoc
//
//	@Summary		List Games
//	@Description	get games
//	@Security		JWT
//	@Tags			games
//	@Accept			json
//	@Produce		json
//	@Param			title		query		string		false	"game search by title"
//	@Param			platforms	query		[]string	false	"game search by platforms, comma-separated"
//	@Param			page		query		int			false	"page number, default is 1"
//	@Success		200			{array}		mappers.PaginationResponse[[]mappers.GameOutputDTO]
//	@Failure		400			{object}	mappers.ErrorResponse
//	@Failure		404			{object}	mappers.PaginationResponse[[]mappers.GameOutputDTO]
//	@Failure		500			{object}	mappers.ErrorResponse
//	@Router			/games [get]
func (h *GameHandler) ListGames(c *fiber.Ctx) error {
	ctx := c.Context()
	pageStr := c.Query("page", "1")          // Default page to "1"
	platformsStr := c.Query("platforms", "") // Default to empty string if not provided
	title := c.Query("title", "")            // Default to empty string if not provided
	if platformsStr == "" && title == "" {
		return c.Status(http.StatusBadRequest).JSON(mappers.ErrorResponse{
			Error:   "At least one of 'platforms' or 'title' query parameters must be provided",
			Details: "Please provide at least one search criterion.",
		})
	}
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1 // Default to page 1 if conversion fails or page is invalid
	}
	platforms := utils.SanitizeArrayStrings(platformsStr)
	log.Printf("Handler: Listing games with page: %d, platforms: '%s', title: '%s'", page, platforms, title)
	games, total, err := h.gameService.ListGames(ctx, page, platforms, title)
	if err != nil {
		log.Printf("Error from GameService: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(mappers.ErrorResponse{
			Error:   "An unexpected error occurred",
			Details: err.Error(),
		})
	}

	if len(games) == 0 {
		return c.Status(http.StatusNotFound).JSON(mappers.PaginationResponse[[]mappers.GameOutputDTO]{
			CommonResponse: mappers.CommonResponse[[]mappers.GameOutputDTO]{
				Data:    []mappers.GameOutputDTO{},
				Message: "No games found matching your criteria",
			},
			Filters: fiber.Map{
				"platforms": platforms,
				"title":     title,
			},
			Page:  page,
			Total: 0,
			Count: 0,
		})
	}

	return c.Status(http.StatusOK).JSON(mappers.PaginationResponse[[]mappers.GameOutputDTO]{
		CommonResponse: mappers.CommonResponse[[]mappers.GameOutputDTO]{
			Data:    mappers.MapGamesModelToOutputDTO(games),
			Message: "Games retrieved successfully",
		},
		Filters: fiber.Map{
			"platforms": platforms,
			"title":     title,
		},
		Page:  page,
		Total: total,
		Count: len(games),
	})
}
