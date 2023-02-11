package rest

import (
	"net/http"
)

type FailIfActivatedHandler struct {
	userConfig UserConfig
	handler    http.Handler
}

func NewFailIfActivatedHandler(userConfig UserConfig, handler http.Handler) *FailIfActivatedHandler {
	return &FailIfActivatedHandler{
		userConfig: userConfig,
		handler:    handler,
	}
}

func (h *FailIfActivatedHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h.userConfig.IsActivated() {
		http.Error(w, "Device is activated", 502)
		return
	}
	h.handler.ServeHTTP(w, r)
}
