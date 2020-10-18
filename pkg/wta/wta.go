package wta

import (
	"context"
	"net/http"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	baseURL = "https://www.wta.org/@@search_tripreport_listing/"
)

var photosRE = regexp.MustCompile(`([0-9]+)`)

// Client .
type Client struct {
	client    *http.Client
	collector *colly.Collector

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
		client:    &http.Client{},
		collector: NewCollector(),
	}
	for _, opt := range opts {
		err := opt(c)
		if err != nil {
			return nil, err
		}
	}

	// Services used for talking to the WTA
	c.Reports = &ReportsService{client: c}
	c.Regions = &RegionsService{client: c}

	return c, nil
}

// NewCollector .
func NewCollector() *colly.Collector {
	c := colly.NewCollector(
		colly.AllowedDomains("wta.org", "www.wta.org"),
	)
	c.SetRequestTimeout(10 * time.Second)
	return c
}

// WithTimeout timeout
func WithTimeout(timeout time.Duration) func(*Client) error {
	return func(client *Client) error {
		client.client.Timeout = timeout
		return nil
	}
}

// WithHTTPClient .
func WithHTTPClient(client *http.Client) func(c *Client) error {
	return func(c *Client) error {
		if client != nil {
			c.client = client
		}
		return nil
	}
}

// WithCollector collector
func WithCollector(collector *colly.Collector) func(*Client) error {
	return func(c *Client) error {
		if collector != nil {
			c.collector = collector
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
