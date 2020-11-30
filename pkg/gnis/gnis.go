package gnis

//go:generate go run ../../cmd/genwith/genwith.go --client --package gnis

import (
	"net/http"
)

const (
	baseURL = "https://geonames.usgs.gov"
)

// Client used to communicate with GNIS
type Client struct {
	client *http.Client

	GeoNames *GeoNamesService
}

func withServices(c *Client) {
	c.GeoNames = &GeoNamesService{client: c}
}
