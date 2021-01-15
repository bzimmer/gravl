package gnis

//go:generate genwith --client --package gnis

import "net/http"

const baseURL = "https://geonames.usgs.gov"

// Client provides access to the GNIS database
type Client struct {
	client *http.Client

	GeoNames *GeoNamesService
}

func withServices() Option {
	return func(c *Client) error {
		c.GeoNames = &GeoNamesService{client: c}
		return nil
	}
}
