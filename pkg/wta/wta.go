package wta

import (
	"context"
	"net/http"
	"time"

	"github.com/bzimmer/transport"
	"github.com/gin-gonic/gin"
)

const (
	baseURL = "https://www.wta.org/@@search_tripreport_listing/"
)

// Client .
type Client struct {
	client *http.Client

	Reports *ReportsService
	Regions *RegionsService
}

type service struct {
	client *Client //nolint:golint,structcheck
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

	// Services used for communicating with the WTA website
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

// WithTransport transport
func WithTransport(transport http.RoundTripper) Option {
	return func(c *Client) error {
		if transport != nil {
			c.client.Transport = transport
		}
		return nil
	}
}

// WithHTTPTracing .
func WithHTTPTracing(debug bool) Option {
	return func(c *Client) error {
		if !debug {
			return nil
		}
		c.client.Transport = &transport.VerboseTransport{
			Transport: c.client.Transport,
		}
		return nil
	}
}

// TripReportsHandler .
func TripReportsHandler(client *Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		reporter := c.Param("reporter")
		reports, err := client.Reports.TripReports(context.Background(), reporter)
		if err != nil {
			c.Abort()
			_ = c.Error(err)
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
