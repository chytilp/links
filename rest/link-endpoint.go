package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"strconv"
	"time"

	"github.com/chytilp/links/datalayer"
	"github.com/chytilp/links/model"
)

// LinkHandler type is type for handling requests to link endpoint.
type LinkHandler struct{}

func (h *LinkHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var err error
	switch r.Method {
	case "GET":
		err = h.handleGet(w, r)
	case "POST":
		err = h.handlePost(w, r)
	case "PUT":
		err = h.handlePut(w, r)
	case "DELETE":
		err = h.handleDelete(w, r)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *LinkHandler) handleGet(w http.ResponseWriter, r *http.Request) error {
	var err error
	urlPath := path.Base(r.URL.Path)
	if urlPath == "link" {
		return h.handleRetrieve(w, r)
	}
	id, err := strconv.Atoi(urlPath)
	if err != nil {
		outErr := fmt.Errorf("Path parameter wrong type, value: %s . Error: %s", urlPath, err)
		prepareResponseFromError(w, outErr, 404)
	}
	links := datalayer.CreateLinks(nil)
	defer links.Close()
	link, err := links.Get(int(id))
	if err != nil {
		outErr := fmt.Errorf("Link with id=%d was not found. Error: %s", id, err)
		prepareResponseFromError(w, outErr, 404)
	} else {
		output, _ := json.Marshal(link)
		prepareResponseFromBytes(w, output, 200)
	}
	return nil
}

func (h *LinkHandler) handleRetrieve(w http.ResponseWriter, r *http.Request) error {
	queryParams := r.URL.Query()
	links := datalayer.CreateLinks(nil)
	defer links.Close()
	foundLinks, err := links.Retrieve(queryParams)
	if err != nil {
		return err
	}
	content := make([]string, len(foundLinks))
	for index, link := range foundLinks {
		bytes, _ := json.Marshal(link)
		content[index] = string(bytes)
	}
	output, err := json.Marshal(content)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(output)
	return nil
}

func (h *LinkHandler) handlePost(w http.ResponseWriter, r *http.Request) error {
	err := h.processSave(w, r)
	if err != nil {
		return err
	}
	return nil
}

func (h *LinkHandler) processSave(w http.ResponseWriter, r *http.Request) error {
	len := r.ContentLength
	body := make([]byte, len)
	r.Body.Read(body)
	var link model.Link
	json.Unmarshal(body, &link)
	links := datalayer.CreateLinks(nil)
	defer links.Close()
	var outLink *model.Link
	outLink, err := links.Save(link)
	if err != nil {
		return err
	}
	idmap := make(map[string]int)
	idmap["id"] = outLink.ID
	output, err := json.Marshal(idmap)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	w.Write(output)
	return nil
}

func (h *LinkHandler) handlePut(w http.ResponseWriter, r *http.Request) error {
	err := h.processSave(w, r)
	if err != nil {
		return err
	}
	return nil
}

func (h *LinkHandler) handleDelete(w http.ResponseWriter, r *http.Request) error {
	id, err := strconv.Atoi(path.Base(r.URL.Path))
	if err != nil {
		return err
	}
	links := datalayer.CreateLinks(nil)
	defer links.Close()
	link, err := links.Get(int(id))
	if err != nil {
		return err
	}
	now := time.Now()
	_, err = links.Delete(int(id), now)
	if err != nil {
		return err
	}
	output, _ := json.Marshal(link)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(output)
	return nil
}
