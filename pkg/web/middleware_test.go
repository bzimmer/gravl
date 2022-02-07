package web_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"

	"github.com/bzimmer/gravl/pkg/web"
)

type tokenhandler struct{}

func (t *tokenhandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte(`{
		"access_token":"99881100332255",
		"token_type":"bearer",
		"expires_in":3600,
		"refresh_token":"ThisIsGood",
		"scope":"user"
	  }`))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func TestLogHandler(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	logger := web.NewLogHandler(&log.Logger)

	mux := http.NewServeMux()
	mux.Handle("/token", logger(new(tokenhandler)))

	w := httptest.NewRecorder()
	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "/token", http.NoBody)
	mux.ServeHTTP(w, req)
	a.Equal(http.StatusOK, w.Code)
}
