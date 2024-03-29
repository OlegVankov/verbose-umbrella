package handler

import (
	"bytes"
	"compress/gzip"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"strings"

	"go.uber.org/zap"

	"github.com/OlegVankov/verbose-umbrella/internal/logger"
)

type gzipWriter struct {
	http.ResponseWriter
	gw *gzip.Writer
}

func newGzipWriter(rw http.ResponseWriter) *gzipWriter {
	return &gzipWriter{rw, gzip.NewWriter(rw)}
}

func (w *gzipWriter) Write(b []byte) (int, error) {
	contentType := strings.Join(w.Header().Values("Content-Type"), ",")
	gzipAccept := strings.Contains(contentType, "application/json") || strings.Contains(contentType, "text/html")
	if gzipAccept {
		w.Header().Set("Content-Encoding", "gzip")
		w.WriteHeader(http.StatusOK)
		return w.gw.Write(b)
	}
	return len(b), nil
}

func (w *gzipWriter) Close() error {
	return w.gw.Close()
}

type gzipReader struct {
	reader io.ReadCloser
	gr     *gzip.Reader
}

func newGzipReader(reader io.ReadCloser) (*gzipReader, error) {
	gr, err := gzip.NewReader(reader)
	if err != nil {
		logger.Log.Info("Error create gzip reader", zap.Error(err))
		return nil, err
	}
	return &gzipReader{
		reader: reader,
		gr:     gr,
	}, nil
}

func (g gzipReader) Read(p []byte) (n int, err error) {
	return g.gr.Read(p)
}

func (g *gzipReader) Close() error {
	if err := g.reader.Close(); err != nil {
		return err
	}
	return g.gr.Close()
}

func compressMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Encoding") == "gzip" {
			body, err := newGzipReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			defer body.Close()
			r.Body = body
		}

		acceptEncoding := strings.Join(r.Header.Values("Accept-Encoding"), ",")
		if !strings.Contains(acceptEncoding, "gzip") {
			h.ServeHTTP(w, r)
			return
		}

		writer := newGzipWriter(w)
		defer writer.Close()
		w = writer

		h.ServeHTTP(w, r)
	})
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
		hex.EncodeToString(h1.Sum(nil))

		// if mac1 != mac2 {
		// 	w.WriteHeader(http.StatusBadRequest)
		// 	return
		// }

		body := io.NopCloser(bytes.NewBuffer(buf))
		r.Body = body
		next.ServeHTTP(w, r)
	})
}
