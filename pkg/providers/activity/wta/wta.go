package wta

//go:generate genwith --client --package wta

import "net/http"

const baseURL = "https://www.wta.org/@@search_tripreport_listing/"

// Client for accessing WTA trip reports
type Client struct {
	client *http.Client

	Reports *ReportsService
	Regions *RegionsService
}

func withServices() Option {
	return func(c *Client) error {
		c.Reports = &ReportsService{client: c}
		c.Regions = &RegionsService{client: c}
		return nil
	}
}
