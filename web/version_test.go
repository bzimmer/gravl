package web_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/bzimmer/gravl/web"
)

func TestVersionHandler(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	mux := http.NewServeMux()
	mux.HandleFunc("/version", web.VersionHandler())

	w := httptest.NewRecorder()
	req, _ := http.NewRequestWithContext(t.Context(), http.MethodGet, "/version", http.NoBody)
	mux.ServeHTTP(w, req)
	a.Equal(http.StatusOK, w.Code)

	var version map[string]any
	decoder := json.NewDecoder(w.Body)
	err := decoder.Decode(&version)
	a.NoError(err)
	a.Equal("development", version["build_version"])
}

// errorWriter forces Write to fail so that json encoding errors are exercised.
type errorWriter struct {
	rec *httptest.ResponseRecorder
}

func (e *errorWriter) Header() http.Header          { return e.rec.Header() }
func (e *errorWriter) WriteHeader(code int)         { e.rec.WriteHeader(code) }
func (e *errorWriter) Write(_ []byte) (int, error) { return 0, errors.New("forced write error") }

func TestVersionHandler_WriteError(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	handler := web.VersionHandler()
	inner := httptest.NewRecorder()
	req, _ := http.NewRequestWithContext(t.Context(), http.MethodGet, "/version", http.NoBody)
	handler(&errorWriter{rec: inner}, req)
	a.Equal(http.StatusInternalServerError, inner.Code)
}
