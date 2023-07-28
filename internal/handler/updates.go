package handler

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/OlegVankov/verbose-umbrella/internal/logger"
	"github.com/OlegVankov/verbose-umbrella/internal/storage"
)

func (h *Handler) updates(w http.ResponseWriter, req *http.Request) {
	var buf bytes.Buffer
	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var metrics []storage.Metrics
	if err := json.Unmarshal(buf.Bytes(), &metrics); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	errs := []string{}
	for _, m := range metrics {
		switch m.MType {
		case "counter":
			delta, err := h.Storage.UpdateCounter(req.Context(), m.ID, *m.Delta)
			if err != nil {
				errs = append(errs, err.Error())
				continue
			}
			m.Delta = &delta
		case "gauge":
			err := h.Storage.UpdateGauge(req.Context(), m.ID, *m.Value)
			if err != nil {
				errs = append(errs, err.Error())
				continue
			}
		}
	}

	// если есть ошибки логируем их
	// если количество ошибок равно количеству метрик выходим с ошибкой )
	if len(errs) != 0 {
		logger.Log.Info(strings.Join(errs, "|"))
		if len(errs) == len(metrics) {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status": "ok"}`))
}

func (h *Handler) checkHash(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mac1 := r.Header.Get("HashSHA256")

		if len(mac1) == 0 || len(h.Key) == 0 {
			next.ServeHTTP(w, r)
			return
		}

		buf, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		h1 := hmac.New(sha256.New, []byte(h.Key))
		h1.Write(buf)
		// mac2 := hex.EncodeToString(h1.Sum(nil))

		// if mac1 != mac2 {
		// 	w.WriteHeader(http.StatusBadRequest)
		// 	return
		// }

		body := io.NopCloser(bytes.NewBuffer(buf))
		r.Body = body
		next.ServeHTTP(w, r)
	})
}
