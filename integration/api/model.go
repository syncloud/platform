package api

type Response struct {
	Success bool    `json:"success"`
	Message string  `json:"message,omitempty"`
	Data    *string `json:"data,omitempty"`
}

type ServiceRestart struct {
	Name string `json:"name"`
}

type DkimKey struct {
	Key string `json:"dkim_key"`
}
