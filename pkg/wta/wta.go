package wta

//go:generate go run ../../dev/genwith.go --package wta

import (
	"context"
	"net/http"
	"time"

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
func RegionsHandler(client *Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		regions, err := client.Regions.Regions(context.Background())
		if err != nil {
			c.Abort()
			_ = c.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"status": "failed"})
			return
		}
		c.IndentedJSON(http.StatusOK, regions)
	}
}
