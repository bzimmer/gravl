package strava

import (
	"net/http"
	"strconv"
	"strings"
)

const (
	// HeaderRateLimit header
	HeaderRateLimit = "X-Ratelimit-Limit"
	// HeaderRateUsage header
	HeaderRateUsage = "X-Ratelimit-Usage"
)

// RateLimit .
// http://developers.strava.com/docs/rate-limits/
type RateLimit struct {
	LimitWindow int `json:"limit_window"`
	LimitDaily  int `json:"limit_daily"`
	UsageWindow int `json:"usage_window"`
	UsageDaily  int `json:"usage_daily"`
}

type rateLimitedTransport struct {
	rateLimit *RateLimit
	transport http.RoundTripper
}

// RoundTrip .
func (r *rateLimitedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if r.rateLimit != nil && r.rateLimit.IsThrottled() {
		return nil, r.rateLimit.newError()
	}

	res, err := r.transport.RoundTrip(req)

	if res != nil {
		rl, err := parseRateLimit(res)
		if err != nil {
			return nil, err
		}
		r.rateLimit = rl
	}

	return res, err
}

// RateLimitError .
type RateLimitError struct {
	RateLimit *RateLimit
}

func (e *RateLimitError) Error() string {
	return "exceeded rate limit"
}

func newRateLimitError(rl *RateLimit) *RateLimitError {
	return &RateLimitError{
		RateLimit: rl,
	}
}

// PercentDaily .
func (r *RateLimit) PercentDaily() int {
	if r.LimitDaily == 0 {
		return 0
	}
	return int(float32(r.UsageDaily) / float32(r.LimitDaily) * 100)
}

// PercentWindow .
func (r *RateLimit) PercentWindow() int {
	if r.LimitWindow == 0 {
		return 0
	}
	return int(float32(r.UsageWindow) / float32(r.LimitWindow) * 100)
}

// IsThrottled .
func (r *RateLimit) IsThrottled() bool {
	return r.PercentDaily() >= 100.0 || r.PercentWindow() >= 100.0
}

func (r *RateLimit) newError() *RateLimitError {
	return newRateLimitError(r)
}

// parseRateLimit parses the headers returned from an API call into
// a RateLimit struct
//
//   HTTP/1.1 200 OK
//   Content-Type: application/json; charset=utf-8
//   Date: Tue, 10 Oct 2020 20:11:01 GMT
//   X-Ratelimit-Limit: 600,30000
//   X-Ratelimit-Usage: 314,27536
func parseRateLimit(res *http.Response) (*RateLimit, error) {
	var rateLimit RateLimit
	if limit := res.Header.Get(HeaderRateLimit); limit != "" {
		limits := strings.Split(limit, ",")
		rateLimit.LimitWindow, _ = strconv.Atoi(limits[0])
		rateLimit.LimitDaily, _ = strconv.Atoi(limits[1])
	}
	if usage := res.Header.Get(HeaderRateUsage); usage != "" {
		usages := strings.Split(usage, ",")
		rateLimit.UsageWindow, _ = strconv.Atoi(usages[0])
		rateLimit.UsageDaily, _ = strconv.Atoi(usages[1])
	}
	return &rateLimit, nil
}

// WithRateLimitThrottling .
func WithRateLimitThrottling() func(client *Client) error {
	return func(client *Client) error {
		transport := client.client.Transport
		if transport == nil {
			transport = http.DefaultTransport
		}
		client.client.Transport = &rateLimitedTransport{
			transport: transport,
			rateLimit: &RateLimit{},
		}
		return nil
	}
}
