package handlers

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/paveloborin/fetch-task/pkg/cache"
)

type Handler struct {
	cache *cache.Storage
}

func NewHandler(cache *cache.Storage) *Handler {
	return &Handler{
		cache: cache,
	}
}

func writeResponse(w http.ResponseWriter, code int, resp interface{}) {
	if resp == nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(code)
		return
	}

	data, err := json.Marshal(resp)
	if err != nil {
		log.Printf("Couldn't encode response %+v to HTTP response body.", resp)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	//	len, _ := w.Write(data)

	len := binary.Size(data)
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len))
	w.WriteHeader(code)
	w.Write(data)
}
