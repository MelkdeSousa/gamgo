package models

type Game struct {
	ID          string
	Title       string
	Description string
	ReleaseDate string
	Platforms   []string
	Rating      int
}
