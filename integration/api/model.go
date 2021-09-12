package api

type Response struct {
	Success bool    `json:"success"`
	Message string  `json:"message,omitempty"`
	Data    *string `json:"data,omitempty"`
}
