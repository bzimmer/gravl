package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"go/format"
	"html/template"
	"io/ioutil"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/rs/zerolog/log"
)

type With struct {
	Auth    bool
	Package string
}

const (
	q = `// Code generated by "genwith.go"; DO NOT EDIT.

	package {{.Package}}

import (
{{ with .Auth}}
	"time"
	"golang.org/x/oauth2"
{{end}}
	"net/http"
	"github.com/bzimmer/httpwares"
)

{{ with .Auth}}
// WithConfig sets the underlying config
func WithConfig(config oauth2.Config) Option {
	return func(c *Client) error {
		c.config = config
		return nil
	}
}

// WithTokenCredentials provides the tokens for an authenticated user
func WithTokenCredentials(accessToken, refreshToken string, expiry time.Time) Option {
	return func(c *Client) error {
		c.token.AccessToken = accessToken
		c.token.RefreshToken = refreshToken
		c.token.Expiry = expiry
		return nil
	}
}

// WithAPICredentials provides the client api credentials for the application
func WithClientCredentials(clientID, clientSecret string) Option {
	return func(c *Client) error {
		c.config.ClientID = clientID
		c.config.ClientSecret = clientSecret
		return nil
	}
}
{{end}}

// WithHTTPTracing enables tracing http calls
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

// WithTransport sets the underlying http client transport
func WithTransport(t http.RoundTripper) Option {
	return func(c *Client) error {
		if t != nil {
			c.client.Transport = t
		}
		return nil
	}
}

// WithHTTPClient sets the underlying http client
func WithHTTPClient(client *http.Client) Option {
	return func(c *Client) error {
		if client != nil {
			c.client = client
		}
		return nil
	}
}`
)

func generate(w With) error {
	tmpl, err := template.New("genwith").Parse(q)
	if err != nil {
		log.Error().Err(err).Msg("parsing template")
		return err
	}

	b := new(bytes.Buffer)
	err = tmpl.Execute(b, w)
	if err != nil {
		log.Error().Err(err).Msg("executing template")
		return err
	}

	// runs the go code formatter
	src, err := format.Source(b.Bytes())
	if err != nil {
		return err
	}

	file := fmt.Sprintf("%s_with.go", w.Package)
	if err := ioutil.WriteFile(file, src, 0644); err != nil {
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
			&cli.StringFlag{
				Name:  "package",
				Value: "",
				Usage: "The name of the package for generation",
			},
		},
		Before: func(c *cli.Context) error {
			if !c.IsSet("package") {
				return errors.New("missing package")
			}
			return nil
		},
		Action: func(c *cli.Context) error {
			w := With{
				Auth:    c.Bool("auth"),
				Package: c.String("package"),
			}
			return generate(w)
		},
	}
	ctx := context.Background()
	if err := app.RunContext(ctx, os.Args); err != nil {
		log.Error().Err(err).Send()
		os.Exit(1)
	}
	os.Exit(0)
}
