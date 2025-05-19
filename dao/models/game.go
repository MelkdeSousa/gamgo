package models

type Game struct {
	ID             string   `json:"id"`
	Title          string   `json:"title"`
	Description    string   `json:"description,omitempty"`
	ReleaseDate    string   `json:"release_date,omitempty"`
	Platforms      []string `json:"platforms"`
	Rating         int      `json:"rating"`
	ExternalID     string   `json:"external_id"`
	ExternalSource string   `json:"external_source"`
	CoverImage     string   `json:"cover_image"`
}
