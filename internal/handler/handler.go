package handler

import (
	storage2 "github.com/OlegVankov/verbose-umbrella/internal/storage"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Handler struct {
	Mux *http.ServeMux
}

func NewHandler() *Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/update/", update)
	return &Handler{
		Mux: mux,
	}
}

var storage = &storage2.MemStorage{
	Gauge:   map[string]storage2.Gauge{},
	Counter: map[string]storage2.Counter{},
}

func update(w http.ResponseWriter, req *http.Request) {
	log.Println(req.URL.Path)
	if req.Method != http.MethodPost {
		w.WriteHeader(http.StatusNotImplemented)
		return
	}
	m := parseURL(req)
	if len(m) != 3 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	switch m[0] {
	case "counter":
		value, err := strconv.ParseInt(m[2], 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		storage.UpdateCounter(m[1], storage2.Counter(value))
	case "gauge":
		value, err := strconv.ParseFloat(m[2], 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		storage.UpdateGauge(m[1], storage2.Gauge(value))
	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "plain/text")
	w.WriteHeader(http.StatusOK)
	log.Println(storage)
	w.Write([]byte(""))
}

func parseURL(r *http.Request) []string {
	path := strings.TrimSpace(r.URL.Path)
	path = strings.TrimPrefix(path, "/update/")

	if strings.HasSuffix(path, "/") {
		sz := len(path) - 1
		path = path[:sz]
	}

	sl := strings.Split(path, "/")

	return sl
}
