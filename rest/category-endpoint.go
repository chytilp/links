package rest

import "net/http"

// CategoryHandler type is type for handling requests to category endpoint.
type CategoryHandler struct{}

func (h *CategoryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {}
