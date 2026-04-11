package web_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"

	"github.com/bzimmer/gravl/web"
)

type tokenhandler struct{}

func (t *tokenhandler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
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
	req, _ := http.NewRequestWithContext(t.Context(), http.MethodGet, "/token", http.NoBody)
	mux.ServeHTTP(w, req)
	a.Equal(http.StatusOK, w.Code)
}

func TestLogHandler_4xx(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	logger := web.NewLogHandler(&log.Logger)
	mux := http.NewServeMux()
	mux.Handle("/bad", logger(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "bad request", http.StatusBadRequest)
	})))

	w := httptest.NewRecorder()
	req, _ := http.NewRequestWithContext(t.Context(), http.MethodGet, "/bad", http.NoBody)
	mux.ServeHTTP(w, req)
	a.Equal(http.StatusBadRequest, w.Code)
}

func TestLogHandler_5xx(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	logger := web.NewLogHandler(&log.Logger)
	mux := http.NewServeMux()
	mux.Handle("/err", logger(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "internal error", http.StatusInternalServerError)
	})))

	w := httptest.NewRecorder()
	req, _ := http.NewRequestWithContext(t.Context(), http.MethodGet, "/err", http.NoBody)
	mux.ServeHTTP(w, req)
	a.Equal(http.StatusInternalServerError, w.Code)
}

func TestLogHandler_DoubleWriteHeader(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	logger := web.NewLogHandler(&log.Logger)
	mux := http.NewServeMux()
	mux.Handle("/double", logger(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.WriteHeader(http.StatusOK) // second call is a no-op
	})))

	w := httptest.NewRecorder()
	req, _ := http.NewRequestWithContext(t.Context(), http.MethodGet, "/double", http.NoBody)
	mux.ServeHTTP(w, req)
	a.Equal(http.StatusOK, w.Code)
}
