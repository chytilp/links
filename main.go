package main

import (
	"net/http"

	"github.com/chytilp/links/rest"
)

func main() {
	server := http.Server{
		Addr: "127.0.0.1:9073",
	}

	http.Handle("/link/", &rest.LinkHandler{})
	http.Handle("/category/", &rest.CategoryHandler{})
	http.Handle("/user/", &rest.UserHandler{})
	http.Handle("/role/", &rest.RoleHandler{})
	server.ListenAndServe()
}
