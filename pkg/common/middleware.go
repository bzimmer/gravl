package common

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/bzimmer/gravl/pkg"
)

// LogMiddleware .
func LogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start)

		msg := "Request"
		if len(c.Errors) > 0 {
			msg = c.Errors.String()
		}

		var entry *zerolog.Event
		switch {
		case c.Writer.Status() >= http.StatusBadRequest && c.Writer.Status() < http.StatusInternalServerError:
			entry = log.Warn()
		case c.Writer.Status() >= http.StatusInternalServerError:
			entry = log.Error()
		default:
			entry = log.Info()
		}

		entry.
			Str("client_ip", c.ClientIP()).
			Dur("elapsed", duration).
			Str("method", c.Request.Method).
			Str("path", c.Request.RequestURI).
			Int("status", c.Writer.Status()).
			Str("referrer", c.Request.Referer()).
			Str("user_agent", c.Request.Header.Get("User-Agent")).
			Msg(msg)
	}
}

// VersionHandler .
func VersionHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.IndentedJSON(http.StatusOK, map[string]string{
			"build_version": pkg.BuildVersion,
		})
	}
}
