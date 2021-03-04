package rest

import "net/http"

// RoleHandler type is type for handling requests to role endpoint.
type RoleHandler struct{}

func (h *RoleHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {}
