package mappers

type CommonResponse[D any] struct {
	Data    D      `json:"data"`
	Message string `json:"message,omitempty"`
}

type AuthResponse struct {
	Token        string `json:"token"`
	ExpirationAt int64  `json:"expiration"` // Unix timestamp for token expiration
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Details string `json:"details,omitempty"`
}

type PaginationResponse[D any] struct {
	CommonResponse[D]
	Page    int         `json:"page"`
	Total   int         `json:"total,omitempty"`
	Count   int         `json:"count"`
	Filters interface{} `json:"filters"`
}
