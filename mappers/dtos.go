package mappers

type CommonResponse[D any] struct {
	Data    D      `json:"data"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Details string `json:"details"`
}

type PaginationResponse[D any] struct {
	CommonResponse[D]
	Page    int         `json:"page"`
	Total   int         `json:"total,omitempty"`
	Count   int         `json:"count"`
	Filters interface{} `json:"filters"`
}
