package model

import "encoding/json"

type ServerResponse struct {
	Result json.RawMessage `json:"result"`
	Status string          `json:"status"`
}
