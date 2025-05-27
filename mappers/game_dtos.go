package mappers

type GameOutputDTO struct {
	Id        string   `json:"id"`
	Title     string   `json:"title"`
	Released  string   `json:"released"`
	Platforms []string `json:"platforms"`
	Rating    float64  `json:"rating"`
}
