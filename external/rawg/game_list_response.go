package rawg

type GameListResponse struct {
	Count         int64       `json:"count"`
	Next          interface{} `json:"next"`
	Previous      interface{} `json:"previous"`
	Results       []Result    `json:"results"`
	UserPlatforms bool        `json:"user_platforms"`
}

type Result struct {
	Slug            string     `json:"slug"`
	Name            string     `json:"name"`
	Platforms       []Platform `json:"platforms"`
	ID              int        `json:"id"`
	Released        string     `json:"released"`
	Rating          float64    `json:"rating"`
	BackgroundImage *string    `json:"background_image"`
	// Playtime         int64             `json:"playtime"`
	// Stores           []Store           `json:"stores"`
	// Tba              bool              `json:"tba"`
	// RatingTop        int64             `json:"rating_top"`
	// Ratings          []Rating          `json:"ratings"`
	// RatingsCount     int64             `json:"ratings_count"`
	// ReviewsTextCount int64             `json:"reviews_text_count"`
	// Added            int64             `json:"added"`
	// AddedByStatus    *AddedByStatus    `json:"added_by_status"`
	// Metacritic       *int64            `json:"metacritic"`
	// SuggestionsCount int64             `json:"suggestions_count"`
	// Updated          time.Time         `json:"updated"`
	// Score            string            `json:"score"`
	// Clip             interface{}       `json:"clip"`
	// Tags             []Tag             `json:"tags"`
	// EsrbRating       interface{}       `json:"esrb_rating"`
	// UserGame         interface{}       `json:"user_game"`
	// ReviewsCount     int64             `json:"reviews_count"`
	// SaturatedColor   string            `json:"saturated_color"`
	// DominantColor    string            `json:"dominant_color"`
	// ShortScreenshots []ShortScreenshot `json:"short_screenshots"`
	// ParentPlatforms  []Platform        `json:"parent_platforms"`
	// Genres           []Genre           `json:"genres"`
	// CommunityRating  *int64            `json:"community_rating,omitempty"`
}

type AddedByStatus struct {
	Yet     *int64 `json:"yet,omitempty"`
	Owned   int64  `json:"owned"`
	Beaten  *int64 `json:"beaten,omitempty"`
	Toplay  *int64 `json:"toplay,omitempty"`
	Dropped *int64 `json:"dropped,omitempty"`
	Playing *int64 `json:"playing,omitempty"`
}

type Genre struct {
	Name string `json:"name"`
	// ID   int64  `json:"id"`
	// Slug string `json:"slug"`
}

type Platform struct {
	Platform Genre `json:"platform"`
}

type Rating struct {
	ID      int64   `json:"id"`
	Title   string  `json:"title"`
	Count   int64   `json:"count"`
	Percent float64 `json:"percent"`
}

type ShortScreenshot struct {
	ID    int64  `json:"id"`
	Image string `json:"image"`
}

type Store struct {
	Store Genre `json:"store"`
}

type Tag struct {
	ID              int64    `json:"id"`
	Name            string   `json:"name"`
	Slug            string   `json:"slug"`
	Language        Language `json:"language"`
	GamesCount      int64    `json:"games_count"`
	ImageBackground string   `json:"image_background"`
}

type Language string

const (
	Eng Language = "eng"
	Rus Language = "rus"
)
