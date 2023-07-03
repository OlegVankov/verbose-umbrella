package handler

import (
	"compress/gzip"
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
	return &gzipWriter{
		ResponseWriter: rw,
		gw:             gzip.NewWriter(rw),
	}
}

func (g *gzipWriter) Write(bytes []byte) (int, error) {
	return g.gw.Write(bytes)
}

func (g *gzipWriter) Close() error {
	return g.gw.Close()
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
		writer := w
		contentEncoding := strings.Join(r.Header.Values("Content-Encoding"), ",")
		acceptEncoding := strings.Join(r.Header.Values("Accept-Encoding"), ",")

		if strings.Contains(contentEncoding, "gzip") {
			body, err := newGzipReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			r.Body = body
			defer body.Close()
		}

		if strings.Contains(acceptEncoding, "gzip") {
			logger.Log.Info("Middleware compress", zap.Strings("Content-Type", writer.Header().Values("Content-Type")))
			w.Header().Set("Content-Encoding", "gzip")
			gz := newGzipWriter(w)
			writer = gz
			defer gz.Close()
		}

		h.ServeHTTP(writer, r)
	})
}
