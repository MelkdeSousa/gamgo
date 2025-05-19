package rawg

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/melkdesousa/gamgo/config"
)

type RawgAPI struct {
	nextPage     string
	previousPage string
}

func NewRawgAPI() *RawgAPI {
	return &RawgAPI{}
}

func (api *RawgAPI) SearchGames(ctx context.Context, title string, page int) (*GameListResponse, error) {
	baseURL := fmt.Sprintf("https://api.rawg.io/api/games?key=%s&search=%s&page=%d", config.MustGetEnv("RAWG_API_KEY"), url.QueryEscape(title), page)
	log.Printf("Searching games with title: %s, page: %d", title, page)
	resp, err := http.Get(baseURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var response GameListResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}
	return &response, nil
}
