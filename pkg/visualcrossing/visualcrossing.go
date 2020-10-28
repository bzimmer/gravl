package visualcrossing

import "net/http"

// https://www.visualcrossing.com/resources/documentation/weather-api/weather-api-documentation/

const (
	visualCrossingURI = "https://weather.visualcrossing.com/VisualCrossingWebServices/rest/services/weatherdata"
	userAgent         = "(github.com/bzimmer/wta/visualcrossing)"
)

// Client .
type Client struct {
	header http.Header
	client *http.Client

	Forecast *ForecastService
}

type service struct {
	client *Client
}

// Option .
type Option func(*Client) error

// NewClient .
func NewClient(opts ...Option) (*Client, error) {
	c := &Client{
		client: &http.Client{},
		header: make(http.Header),
	}
	// set now, possibly overwritten with options
	c.header.Set("User-Agent", userAgent)
	c.header.Set("Accept", "application/json")
	for _, opt := range opts {
		err := opt(c)
		if err != nil {
			return nil, err
		}
	}

	// Services used for communicating with VisualCrossing
	c.Forecast = &ForecastService{client: c}

	return c, nil
}
