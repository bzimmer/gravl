package strava

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/urfave/cli/v2"
	"golang.org/x/oauth2"

	"github.com/bzimmer/gravl"
)

// Window holds a rate limit figure for Strava's 15-minute and daily windows.
type Window struct {
	FifteenMinute int `json:"15m"`
	Daily         int `json:"daily"`
}

// Limits pairs a limit with its current usage.
type Limits struct {
	Limit Window `json:"limit"`
	Usage Window `json:"usage"`
}

// RateLimits holds the rate limits Strava reports on every response: overall, and for read
// requests when the response carries them.
type RateLimits struct {
	Overall Limits  `json:"overall"`
	Read    *Limits `json:"read,omitempty"`
}

// RateLimitTransport records the rate limit headers of every response passing through it.
type RateLimitTransport struct {
	Transport http.RoundTripper

	mu     sync.Mutex
	limits *RateLimits
}

// RoundTrip implements http.RoundTripper.
func (t *RateLimitTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	next := t.Transport
	if next == nil {
		next = http.DefaultTransport
	}
	res, err := next.RoundTrip(req)
	if err != nil {
		return res, err
	}
	if limits, ok := parseRateLimits(res.Header); ok {
		t.mu.Lock()
		t.limits = limits
		t.mu.Unlock()
	}
	return res, nil
}

// Reset forgets the recorded rate limits.
func (t *RateLimitTransport) Reset() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.limits = nil
}

// RateLimits returns the most recently recorded rate limits, or nil if no response carried them.
func (t *RateLimitTransport) RateLimits() *RateLimits {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.limits
}

// rateLimits records the headers of every request made by the Strava client created in Before.
var rateLimits = &RateLimitTransport{} //nolint:gochecknoglobals // shared by Before and ratelimits

// RateLimitRecorder returns the transport recording the Strava client's rate limit headers.
func RateLimitRecorder() *RateLimitTransport {
	return rateLimits
}

// withRateLimitRecorder returns a context whose oauth2 client sends requests through the
// rate limit recorder, which is how strava.WithAutoRefresh builds its transport.
func withRateLimitRecorder(ctx context.Context) context.Context {
	return context.WithValue(ctx, oauth2.HTTPClient, &http.Client{Transport: rateLimits})
}

func parseWindow(value string) (Window, bool) {
	fifteen, daily, found := strings.Cut(value, ",")
	if !found {
		return Window{}, false
	}
	f, err := strconv.Atoi(strings.TrimSpace(fifteen))
	if err != nil {
		return Window{}, false
	}
	d, err := strconv.Atoi(strings.TrimSpace(daily))
	if err != nil {
		return Window{}, false
	}
	return Window{FifteenMinute: f, Daily: d}, true
}

func parseLimits(header http.Header, limit, usage string) (*Limits, bool) {
	l, ok := parseWindow(header.Get(limit))
	if !ok {
		return nil, false
	}
	u, ok := parseWindow(header.Get(usage))
	if !ok {
		return nil, false
	}
	return &Limits{Limit: l, Usage: u}, true
}

func parseRateLimits(header http.Header) (*RateLimits, bool) {
	overall, ok := parseLimits(header, "X-RateLimit-Limit", "X-RateLimit-Usage")
	if !ok {
		return nil, false
	}
	limits := &RateLimits{Overall: *overall}
	if read, found := parseLimits(header, "X-ReadRateLimit-Limit", "X-ReadRateLimit-Usage"); found {
		limits.Read = read
	}
	return limits, true
}

func ratelimitsCommand() *cli.Command {
	return &cli.Command{
		Name:  "ratelimits",
		Usage: "Query the application's rate limits and current usage",
		Description: "Make one request (the authenticated athlete) and report the rate limits and usage " +
			"Strava returns in its response headers, for the 15-minute and daily windows",
		Action: func(c *cli.Context) error {
			ctx, cancel := context.WithTimeout(c.Context, c.Duration("timeout"))
			defer cancel()
			t := time.Now()
			// report only what this request's response carried, never an earlier one's
			rateLimits.Reset()
			if _, err := gravl.Runtime(c).Strava.Athlete.Athlete(ctx); err != nil {
				return err
			}
			limits := rateLimits.RateLimits()
			if limits == nil {
				return errors.New("no rate limit headers in the response")
			}
			met := gravl.Runtime(c).Metrics
			met.IncrCounter([]string{Provider, c.Command.Name}, 1)
			met.AddSample([]string{Provider, c.Command.Name}, float32(time.Since(t).Seconds()))
			return gravl.Runtime(c).Encoder.Encode(limits)
		},
	}
}
