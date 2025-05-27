package handlers

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/melkdesousa/gamgo/config"
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

func (h *AuthHandler) login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}
	account, err := h.accountService.GetAccount(req.Email, req.Password)
	if err != nil {
		log.Printf("Error getting account: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	if account == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
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
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not login"})
	}
	return c.JSON(fiber.Map{"token": signedToken, "expiration": claims["exp"]})
}
