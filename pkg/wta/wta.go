package wta

//go:generate go run ../../cmd/genwith/genwith.go --client --package wta

import (
	"net/http"
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
