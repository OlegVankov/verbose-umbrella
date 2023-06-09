package handler

import (
	"fmt"
	storage2 "github.com/OlegVankov/verbose-umbrella/internal/storage"
	"github.com/go-chi/chi/v5"
	"html/template"
	"net/http"
	"strconv"
)

var storage = &storage2.MemStorage{
	Gauge:   map[string]storage2.Gauge{},
	Counter: map[string]storage2.Counter{},
}

func Main(w http.ResponseWriter, req *http.Request) {
	ts, err := template.ParseFiles("./html/home.page.tmpl")
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}

	err = ts.Execute(w, storage)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}

}

func Update(w http.ResponseWriter, req *http.Request) {
	typeMetric := chi.URLParam(req, "type")
	name := chi.URLParam(req, "name")

	switch typeMetric {
	case "counter":
		value, err := strconv.ParseInt(chi.URLParam(req, "value"), 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		storage.UpdateCounter(name, storage2.Counter(value))
	case "gauge":
		value, err := strconv.ParseFloat(chi.URLParam(req, "value"), 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		storage.UpdateGauge(name, storage2.Gauge(value))
	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "plain/text")
	w.WriteHeader(http.StatusOK)
}

func ValueGauge(w http.ResponseWriter, req *http.Request) {
	name := chi.URLParam(req, "name")
	value, ok := storage.GetGauge(name)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "plain/text")
	w.WriteHeader(http.StatusOK)

	val := float64(value)
	w.Write([]byte(fmt.Sprint(val)))
}

func ValueCounter(w http.ResponseWriter, req *http.Request) {
	name := chi.URLParam(req, "name")
	value, ok := storage.GetCounter(name)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "plain/text")
	w.WriteHeader(http.StatusOK)

	val := int64(value)
	w.Write([]byte(fmt.Sprint(val)))
}
