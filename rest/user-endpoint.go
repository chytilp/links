package rest

import "net/http"

// UserHandler type is type for handling requests to user endpoint.
type UserHandler struct{}

func (h *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {}
