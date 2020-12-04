package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"text/template"

	"github.com/urfave/cli/v2"

	"github.com/rs/zerolog/log"
)

type with struct {
	Do       bool
	Auth     bool
	Client   bool
	Endpoint bool
	Flags    string
	Package  string
}

const (
	q = `// Code generated by "genwith.go {{.Flags}}"; DO NOT EDIT.

	package {{.Package}}

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/bzimmer/httpwares"
	"golang.org/x/oauth2"
	"io"
	"net/http"
	"time"
)

{{if .Client}}
type service struct {
	client *Client //nolint:golint,structcheck
}

// Option provides a configuration mechanism for a Client
type Option func(*Client) error

// NewClient creates a new client and applies all provided Options
func NewClient(opts ...Option) (*Client, error) {
	c := &Client{
		client: &http.Client{},
	{{- if .Auth}}
		token:  oauth2.Token{},
		config: oauth2.Config{
	{{- if .Endpoint}}
			Endpoint: Endpoint,
	{{end}}
		},
	{{end}}
	}
	for _, opt := range opts {
		err := opt(c)
		if err != nil {
			return nil, err
		}
	}
	withServices(c)
	return c, nil
}
{{end}}

{{if .Auth}}
// WithConfig sets the underlying oauth2.Config.
func WithConfig(config oauth2.Config) Option {
	return func(c *Client) error {
		c.config = config
		return nil
	}
}

// WithToken sets the underlying oauth2.Token.
func WithToken(token oauth2.Token) Option {
	return func(c *Client) error {
		c.token = token
		return nil
	}
}

// WithTokenCredentials provides the tokens for an authenticated user.
func WithTokenCredentials(accessToken, refreshToken string, expiry time.Time) Option {
	return func(c *Client) error {
		c.token.AccessToken = accessToken
		c.token.RefreshToken = refreshToken
		c.token.Expiry = expiry
		return nil
	}
}

// WithAPICredentials provides the client api credentials for the application.
func WithClientCredentials(clientID, clientSecret string) Option {
	return func(c *Client) error {
		c.config.ClientID = clientID
		c.config.ClientSecret = clientSecret
		return nil
	}
}

{{- if .Endpoint}}
// WithAutoRefresh refreshes access tokens automatically.
// The order of this option matters because it is dependent on the client's
// config and token. Use this option after With*Credentials.
func WithAutoRefresh(ctx context.Context) Option {
	return func(c *Client) error {
		c.client = c.config.Client(ctx, &c.token)
		return nil
	}
}
{{end}}
{{end}}

// WithHTTPTracing enables tracing http calls.
func WithHTTPTracing(debug bool) Option {
	return func(c *Client) error {
		if !debug {
			return nil
		}
		c.client.Transport = &httpwares.VerboseTransport{
			Transport: c.client.Transport,
		}
		return nil
	}
}

// WithTransport sets the underlying http client transport.
func WithTransport(t http.RoundTripper) Option {
	return func(c *Client) error {
		if t == nil {
			return errors.New("nil transport")
		}
		c.client.Transport = t
		return nil
	}
}

// WithHTTPClient sets the underlying http client.
func WithHTTPClient(client *http.Client) Option {
	return func(c *Client) error {
		if client == nil {
			return errors.New("nil client")
		}
		c.client = client
		return nil
	}
}

{{if .Do}}
// do executes the http request and populates v with the result.
func (c *Client) do(req *http.Request, v interface{}) error {
	ctx := req.Context()
	res, err := c.client.Do(req)
	if err != nil {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			return err
		}
	}
	defer res.Body.Close()

	httpError := res.StatusCode >= http.StatusBadRequest

	var obj interface{}
	if httpError {
		obj = &Fault{}
	} else {
		obj = v
	}

	if obj != nil {
		err := json.NewDecoder(res.Body).Decode(obj)
		if err == io.EOF {
			err = nil // ignore EOF errors caused by empty response body
		}
		if httpError {
			return obj.(error)
		}
		return err
	}

	return nil
}
{{end}}`

	p = `// Code generated by "genwith.go {{.Flags}}"; DO NOT EDIT.

package {{.Package}}

import (
	"context"
	"errors"

	"github.com/rs/zerolog/log"
)

// Pagination provides guidance on how to paginate through resources
type Pagination struct {
	// Total of resources to query
	Total int
	// Start at this page
	Start int
	// Count of the number of resources to query per page
	Count int
}

// paginator paginates through results
type paginator interface {
	// page returns the default page size
	page() int
	// count of the number of resources queried
	count() int
	// do the querying
	do(ctx context.Context, start, count int) (int, error)
}

func paginate(ctx context.Context, paginator paginator, spec Pagination) error {
	var (
		start = spec.Start
		count = spec.Count
		total = spec.Total
	)
	log.Debug().
		Str("prepared", "pre").
		Int("start", start).
		Int("count", count).
		Int("total", total).
		Msg("paginate")
	if total < 0 {
		return errors.New("total less than zero")
	}
	if start <= 0 {
		start = 1
	}
	if count <= 0 {
		count = paginator.page()
	}
	if total > 0 && total <= count {
		count = total
	}
	// if requesting only one page of data then optimize
	if start <= 1 && total < paginator.page() {
		count = total
	}
	log.Debug().
		Str("prepared", "post").
		Int("start", start).
		Int("count", count).
		Int("total", total).
		Msg("paginate")
	return do(ctx, paginator, total, start, count)
}

func do(ctx context.Context, paginator paginator, total, start, count int) error {
	for {
		all := paginator.count()
		log.Debug().
			Int("all", all).
			Int("start", start).
			Int("count", count).
			Int("total", total).
			Msg("do")
		n, err := paginator.do(ctx, start, count)
		if err != nil {
			return err
		}
		all = paginator.count()
		log.Debug().
			Int("n", n).
			Int("all", all).
			Int("start", start).
			Int("count", count).
			Int("total", total).
			Msg("done")
		// Strava documentation says receiving fewer than requested results is a
		// possible scenario so break only if 0 results were returned or we have
		// enough to fulfill the request
		if n == 0 || all >= total {
			break
		}
		start++
	}
	return nil
}`
)

func format(file string) error {
	cmd := exec.Command("gofmt", "-w", "-s", file)
	if err := cmd.Run(); err != nil {
		return err
	}
	cmd = exec.Command("goimports", "-w", file)
	return cmd.Run()
}

func generate(w with, file, tmpl string) error {
	t, err := template.New("genwith").Parse(tmpl)
	if err != nil {
		log.Error().Err(err).Msg("parsing template")
		return err
	}

	src := new(bytes.Buffer)
	err = t.Execute(src, w)
	if err != nil {
		log.Error().Err(err).Msg("executing template")
		return err
	}

	if err := ioutil.WriteFile(file, src.Bytes(), 0600); err != nil {
		return err
	}
	return nil
}

func main() {
	app := &cli.App{
		Name:     "genwith",
		HelpName: "genwith",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "auth",
				Value: false,
				Usage: "Include auth-related options",
			},
			&cli.BoolFlag{
				Name:  "do",
				Value: false,
				Usage: "Include client.do function",
			},
			&cli.BoolFlag{
				Name:  "client",
				Value: false,
				Usage: "Include NewClient & options",
			},
			&cli.BoolFlag{
				Name:  "endpoint",
				Value: false,
				Usage: "Include oauth2.Endpoint var in config instantiation (--auth must also be enabled)",
			},
			&cli.BoolFlag{
				Name:  "pagination",
				Value: false,
				Usage: "Include pagination framework",
			},
			&cli.StringFlag{
				Name:     "package",
				Value:    "",
				Required: true,
				Usage:    "The name of the package for generation",
			},
		},
		Before: func(c *cli.Context) error {
			if c.Bool("endpoint") && !c.Bool("auth") {
				return errors.New("`endpoint` enabled without `auth`")
			}
			return nil
		},
		Action: func(c *cli.Context) error {
			w := with{
				Do:       c.Bool("do"),
				Auth:     c.Bool("auth"),
				Client:   c.Bool("client"),
				Endpoint: c.Bool("endpoint"),
				Flags:    strings.Join(os.Args[1:], " "),
				Package:  c.String("package")}
			templates := [][]string{
				{fmt.Sprintf("%s_with.go", c.String("package")), q}}
			if c.Bool("pagination") {
				templates = append(templates, []string{"pagination.go", p})
			}
			for _, x := range templates {
				file, template := x[0], x[1]
				if err := generate(w, file, template); err != nil {
					return err
				}
				if err := format(file); err != nil {
					return err
				}
			}
			return nil
		},
	}
	ctx := context.Background()
	if err := app.RunContext(ctx, os.Args); err != nil {
		log.Error().Err(err).Send()
		os.Exit(1)
	}
	os.Exit(0)
}
