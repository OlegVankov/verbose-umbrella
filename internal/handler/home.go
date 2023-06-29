package handler

import (
	"html/template"
	"net/http"

	"github.com/OlegVankov/verbose-umbrella/internal/logger"
	"go.uber.org/zap"
)

func (h *Handler) home(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	ts := template.Must(template.ParseFiles("./html/home.page.tmpl"))
	err := ts.Execute(w, h.storage)
	if err != nil {
		logger.Log.Info("execute template", zap.Error(err))
		w.WriteHeader(http.StatusBadGateway)
		return
	}
}
