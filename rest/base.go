package rest

import (
	"encoding/json"
	"net/http"
)

func prepareResponseFromMap(w http.ResponseWriter, content map[string]string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	output, _ := json.Marshal(content)
	w.Write(output)
}

func prepareResponseFromBytes(w http.ResponseWriter, content []byte, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(content)
}

func prepareResponseFromError(w http.ResponseWriter, err error, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	output, _ := json.Marshal(err)
	w.Write(output)
}
