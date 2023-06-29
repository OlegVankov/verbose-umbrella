package handler

import (
	"compress/gzip"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strings"

	"github.com/OlegVankov/verbose-umbrella/internal/logger"
)

type gzipWriter struct {
	rw http.ResponseWriter
	gw *gzip.Writer
}

func newGzipWriter(rw http.ResponseWriter) *gzipWriter {
	return &gzipWriter{
		rw: rw,
		gw: gzip.NewWriter(rw),
	}
}

func (g *gzipWriter) Header() http.Header {
	return g.rw.Header()
}

func (g *gzipWriter) Write(bytes []byte) (int, error) {
	return g.gw.Write(bytes)
}

func (g *gzipWriter) WriteHeader(statusCode int) {
	g.rw.WriteHeader(statusCode)
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
		nw := w
		contentType := strings.Join(r.Header.Values("Content-Type"), ",")
		contentEncoding := strings.Join(r.Header.Values("Content-Encoding"), ",")
		acceptEncoding := strings.Join(r.Header.Values("Accept-Encoding"), ",")
		accept := strings.Join(r.Header.Values("Accept"), ",")

		if strings.Contains(acceptEncoding, "gzip") &&
			(strings.Contains(contentType, "application/json") || strings.Contains(contentType, "text/html") || strings.Contains(accept, "html/text")) {
			logger.Log.Info("gzip writer", zap.String("Accept-Encoding", acceptEncoding), zap.String("Content-Type", contentType))
			w.Header().Set("Content-Encoding", "gzip")
			gw := newGzipWriter(w)
			defer gw.Close()
			nw = gw
		}

		if strings.Contains(contentEncoding, "gzip") {
			logger.Log.Info("gzip reader", zap.String("Content-Encoding", contentEncoding))
			body, err := newGzipReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			r.Body = body
			defer body.Close()
		}

		h.ServeHTTP(nw, r)
	})
}
