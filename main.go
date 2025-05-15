package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/joho/godotenv"
	"github.com/melkdesousa/gamgo/external/rawg"

	"github.com/gofiber/fiber/v2"
)

func init() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found or error loading .env file")
	}
}

func main() {
	app := fiber.New()

	app.Get("/games/search", func(c *fiber.Ctx) error {
		title := c.Query("title")

		if title == "" {
			return c.Status(400).JSON(fiber.Map{
				"error": "Title query parameter is required",
			})
		}

		// Make HTTP request to search for games
		resp, err := http.Get(fmt.Sprintf("https://api.rawg.io/api/games?key=%s&title=%s&search_exact=true", os.Getenv("RAWG_API_KEY"), url.QueryEscape(title)))
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "Failed to fetch games: " + err.Error(),
			})
		}
		defer resp.Body.Close()

		// Read response body
		body, err := DeserializerJSON[rawg.GameListResponse](resp.Body)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "Failed to read response: " + err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"id":        body.Results[0].ID,
			"name":      body.Results[0].Name,
			"slug":      body.Results[0].Slug,
			"platforms": body.Results[0].Platforms,
		})
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
