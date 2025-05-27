package handlers

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/melkdesousa/gamgo/config"
	"github.com/melkdesousa/gamgo/mappers"
	"github.com/melkdesousa/gamgo/services"
)

var jwtSecret = []byte("your-secret-key") // Use env var in production

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthHandler struct {
	app            *fiber.App
	accountService *services.AccountService // Assuming you have an AccountService for user management
}

func NewAuthHandler(
	app *fiber.App,
	accountService *services.AccountService,
) {
	handler := &AuthHandler{
		app:            app,
		accountService: accountService,
	}
	app.Post("/auth/login", handler.login)
}

// Auth godoc
//
//	@Summary		User Login
//	@Description	Authenticate user and return JWT token
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			login	body		LoginRequest	true	"User login credentials"
//	@Success		200		{object}	mappers.AuthResponse
//	@Failure		400		{object}	mappers.ErrorResponse
//	@Failure		401		{object}	mappers.ErrorResponse
//	@Failure		500		{object}	mappers.ErrorResponse
//	@Router			/auth/login [post]
func (h *AuthHandler) login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(mappers.ErrorResponse{Error: "Invalid request"})
	}
	account, err := h.accountService.GetAccount(req.Email, req.Password)
	if err != nil {
		log.Printf("Error getting account: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(mappers.ErrorResponse{Error: err.Error()})
	}
	if account == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(mappers.ErrorResponse{Error: "Invalid credentials"})
	}
	claims := jwt.MapClaims{
		"email": req.Email,
		"exp":   time.Now().Add(time.Minute * 30).Unix(),
		"iat":   time.Now().Unix(),
		"sub":   account.ID.String(),
		"role":  "user",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(config.MustGetEnv("JWT_SECRET")))
	if err != nil {
		log.Printf("Error signing token: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(mappers.ErrorResponse{Error: "Could not login"})
	}
	return c.JSON(mappers.AuthResponse{Token: signedToken, Expiration: claims["exp"].(int64) - claims["iat"].(int64)}) // Return token and expiration time
}
