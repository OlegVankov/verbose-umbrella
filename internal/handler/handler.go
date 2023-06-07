package handler

import (
	storage2 "github.com/OlegVankov/verbose-umbrella/internal/storage"
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
	Gauge:   map[string]float64{},
	Counter: map[string]int64{},
}

func update(w http.ResponseWriter, req *http.Request) {
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
		storage.Counter[m[1]] += value
	case "gauge":
		value, err := strconv.ParseFloat(m[2], 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		storage.Gauge[m[1]] += value
	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	value, err := strconv.ParseInt(m[2], 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	storage.Counter[m[1]] += value
	w.Header().Set("Content-Type", "plain/text")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(""))
}

func parseURL(r *http.Request) []string {
	path := r.URL.Path
	path = strings.TrimPrefix(strings.TrimSpace(path), "/update/")
	if strings.HasPrefix(path, "/") {
		path = path[1:]
	}
	if strings.HasSuffix(path, "/") {
		sz := len(path) - 1
		path = path[:sz]
	}
	sl := strings.Split(path, "/")
	return sl
}
