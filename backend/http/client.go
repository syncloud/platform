package http

import (
	"net/http"
)

type Client interface {
	Post(url, bodyType string, body interface{}) (*http.Response, error)
}
