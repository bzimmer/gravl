package wta

//go:generate go run ../../../cmd/genwith/genwith.go --client --noservicese --package wta

import "net/http"

const baseURL = "https://www.wta.org/@@search_tripreport_listing/"

// Client for accessing WTA trip reports
type Client struct {
	client *http.Client

	Reports *ReportsService
	Regions *RegionsService
}

func withServices(c *Client) {
	c.Reports = &ReportsService{client: c}
	c.Regions = &RegionsService{client: c}
}
