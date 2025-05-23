package web_test

import (
	"encoding/json"
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
