package models

import "time"

type Game struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	// Description    string    `json:"description,omitempty"`
	ReleaseDate    time.Time `json:"releaseDate,omitempty"`
	Platforms      []string  `json:"platforms"`
	Rating         int       `json:"rating"`
	ExternalID     string    `json:"externalId"`
	ExternalSource string    `json:"externalSource"`
	CoverImage     string    `json:"coverImage"`
}
