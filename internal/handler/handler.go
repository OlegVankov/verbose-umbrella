package handler

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/OlegVankov/verbose-umbrella/internal/storage"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	Router  *chi.Mux
	storage storage.Storage
}

func NewHandler() Handler {
	return Handler{
		Router:  chi.NewRouter(),
		storage: storage.NewStorage(),
	}
}

func (h *Handler) SetRoute() {
	h.Router.Get("/", h.home)
	h.Router.Route("/value", func(r chi.Router) {
		r.Get("/gauge/{name}", h.valueGauge)
		r.Get("/counter/{name}", h.valueCounter)
	})
	h.Router.Route("/update", func(r chi.Router) {
		r.Post("/{type}/{name}/{value}", h.update)
	})
}

func (h *Handler) home(w http.ResponseWriter, req *http.Request) {
	ts, err := template.ParseFiles("./html/home.page.tmpl")
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}

	err = ts.Execute(w, h.storage)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
}

func (h *Handler) update(w http.ResponseWriter, req *http.Request) {
	typeMetric := chi.URLParam(req, "type")
	name := chi.URLParam(req, "name")

	switch typeMetric {
	case "counter":
		value, err := strconv.ParseInt(chi.URLParam(req, "value"), 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		h.storage.UpdateCounter(name, value)
	case "gauge":
		value, err := strconv.ParseFloat(chi.URLParam(req, "value"), 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		h.storage.UpdateGauge(name, value)
	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "plain/text")
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) valueGauge(w http.ResponseWriter, req *http.Request) {
	name := chi.URLParam(req, "name")
	value, ok := h.storage.GetGauge(name)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "plain/text")
	w.WriteHeader(http.StatusOK)

	val := float64(value)
	w.Write([]byte(fmt.Sprint(val)))
}

func (h *Handler) valueCounter(w http.ResponseWriter, req *http.Request) {
	name := chi.URLParam(req, "name")
	value, ok := h.storage.GetCounter(name)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "plain/text")
	w.WriteHeader(http.StatusOK)

	val := int64(value)
	w.Write([]byte(fmt.Sprint(val)))
}
