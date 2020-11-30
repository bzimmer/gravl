package wta

//go:generate go run ../../cmd/genwith/genwith.go --client --package wta

import (
	"context"
	"net/http"

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

func withServices(c *Client) {
	c.Reports = &ReportsService{client: c}
	c.Regions = &RegionsService{client: c}
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
