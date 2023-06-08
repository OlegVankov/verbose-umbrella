package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_update(t *testing.T) {
	tests := []struct {
		name         string
		method       string
		target       string
		expectedCode int
	}{
		{name: "test #1", method: http.MethodGet, target: "/update/gauge/alloc/100", expectedCode: http.StatusNotImplemented},
		{name: "test #2", method: http.MethodPut, target: "/update/gauge/alloc/100", expectedCode: http.StatusNotImplemented},
		{name: "test #3", method: http.MethodDelete, target: "/update/gauge/alloc/100", expectedCode: http.StatusNotImplemented},
		{name: "test #4", method: http.MethodPost, target: "/update/gauge/alloc/100", expectedCode: http.StatusOK},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(tt.method, tt.target, nil)
			w := httptest.NewRecorder()
			update(w, r)
			assert.Equal(t, w.Code, tt.expectedCode, "Код ответа не совпадает с ожидаемым")
		})
	}
}

func Test_parseURL(t *testing.T) {
	type args struct {
		r *http.Request
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "test #1",
			args: struct{ r *http.Request }{r: httptest.NewRequest(http.MethodPost, "/update/gauge/metric/100", nil)},
			want: []string{"gauge", "metric", "100"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, parseURL(tt.args.r), "parseURL(%v)", tt.args.r)
		})
	}
}
