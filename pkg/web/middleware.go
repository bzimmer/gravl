package web

import (
	"net/http"
	"time"

	"github.com/rs/zerolog"
)

type respwriter struct {
	http.ResponseWriter
	wrote  bool
	status int
}

func (w *respwriter) Status() int {
	return w.status
}

func (w *respwriter) Write(p []byte) (n int, err error) {
	if !w.wrote {
		w.WriteHeader(http.StatusOK)
	}
	return w.ResponseWriter.Write(p)
}

func (w *respwriter) WriteHeader(code int) {
	w.ResponseWriter.WriteHeader(code)
	if w.wrote {
		return
	}
	w.wrote = true
	w.status = code
}

// NewLogHandler creates an instance of an http handler used for logging
func NewLogHandler(log zerolog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sw := &respwriter{ResponseWriter: w}
			start := time.Now()
			next.ServeHTTP(sw, r)
			duration := time.Since(start)

			var entry *zerolog.Event
			switch {
			case sw.status >= http.StatusBadRequest && sw.status < http.StatusInternalServerError:
				entry = log.Warn()
			case sw.status >= http.StatusInternalServerError:
				entry = log.Error()
			default:
				entry = log.Info()
			}

			entry.
				Str("client_ip", r.RemoteAddr).
				Dur("elapsed", duration).
				Str("method", r.Method).
				Str("path", r.URL.String()).
				Int("status", sw.status).
				Str("user_agent", r.Header.Get("User-Agent")).
				Msg("request")
		})
	}
}
