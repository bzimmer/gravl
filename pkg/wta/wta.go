package wta

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	baseURL = "https://www.wta.org/@@search_tripreport_listing/"
)

// Client .
type Client struct {
	// collector *colly.Collector
	client *http.Client

	Reports *ReportsService
	Regions *RegionsService
}

type service struct {
	client *Client
}

// Option .
type Option func(*Client) error

// NewClient .
func NewClient(opts ...Option) (*Client, error) {
	c := &Client{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
	for _, opt := range opts {
		err := opt(c)
		if err != nil {
			return nil, err
		}
	}

	// Services used for talking to the WTA website
	c.Reports = &ReportsService{client: c}
	c.Regions = &RegionsService{client: c}

	return c, nil
}

// WithHTTPClient client
func WithHTTPClient(client *http.Client) Option {
	return func(c *Client) error {
		if client != nil {
			c.client = client
		}
		return nil
	}
}

// WithTimeout timeout
func WithTimeout(timeout time.Duration) Option {
	return func(c *Client) error {
		c.client.Timeout = timeout
		return nil
	}
}

// WithTransport transport
func WithTransport(transport http.RoundTripper) Option {
	return func(c *Client) error {
		if transport != nil {
			c.client.Transport = transport
		}
		return nil
	}
}

// ---

// NewRouter .
func NewRouter(client *Client) *gin.Engine {
	r := gin.New()
	r.Use(LogMiddleware(), gin.Recovery())
	r.GET("/version/", VersionHandler())
	r.GET("/regions/", RegionsHandler())
	r.GET("/reports/", TripReportsHandler(client))
	r.GET("/reports/:reporter", TripReportsHandler(client))
	return r
}

// TripReportsHandler .
func TripReportsHandler(client *Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		reporter := c.Param("reporter")
		reports, err := client.Reports.TripReports(context.Background(), reporter)
		if err != nil {
			c.Abort()
			c.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"status": "failed"})
			return
		}
		c.IndentedJSON(http.StatusOK, &TripReports{Reporter: reporter, Reports: reports})
	}
}

// RegionsHandler .
func RegionsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.IndentedJSON(http.StatusOK, Regions)
	}
}

// VersionHandler .
func VersionHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.IndentedJSON(http.StatusOK, map[string]string{
			"build_version": BuildVersion,
		})
	}
}

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
			{
				entry = log.Warn()
			}
		case c.Writer.Status() >= http.StatusInternalServerError:
			{
				entry = log.Error()
			}
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
