package handler

import (
	"html/template"
	"net/http"

	"go.uber.org/zap"

	"github.com/OlegVankov/verbose-umbrella/internal/logger"
)

func (h *Handler) home(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	ts := template.Must(template.ParseFiles("./html/home.page.tmpl"))
	err := ts.Execute(w, h.Storage)
	if err != nil {
		logger.Log.Info("execute template", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
