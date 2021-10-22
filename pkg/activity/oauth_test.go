package activity_test

import (
	"context"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"
	"golang.org/x/sync/errgroup"

	"github.com/bzimmer/gravl/pkg/activity"
	"github.com/bzimmer/gravl/pkg/internal"
)

func command(cfg *activity.OAuthConfig) func(t *testing.T, baseURL string) *cli.Command {
	return func(t *testing.T, baseURL string) *cli.Command {
		return activity.OAuthCommand(cfg)
	}
}

func TestOAuth(t *testing.T) {
	a := assert.New(t)
	tests := []*internal.Harness{
		{
			Name: "success",
			Args: []string{"test", "oauth"},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			/*
				This test works by adding a buffered channel to the config which is then
				selected against in a goroutine to obtain the URL of the oauth callback
				http server.

				A simple request is made to the oauth callback http server but the rest
				of the flow is ignored (it's tested separately).

				The completion of the request to the http server (success or failure)
				ends the goroutine at which time the context is canceled and the http
				server started within the cli.Command will shutdown (if all goes well!).
			*/
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			started := make(chan *url.URL, 1)
			defer close(started)
			cfg := &activity.OAuthConfig{
				Provider: "foobar",
				Started:  started,
				Scopes:   []string{"one", "two", "three"},
			}
			grp, ctx := errgroup.WithContext(ctx)
			grp.Go(func() error {
				defer cancel()
				select {
				case <-ctx.Done():
					return ctx.Err()
				case u := <-started:
					req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
					if err != nil {
						return err
					}
					client := &http.Client{
						CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
							return http.ErrUseLastResponse
						},
					}
					res, err := client.Do(req)
					if err != nil {
						return err
					}
					defer res.Body.Close()
					a.Equal(http.StatusTemporaryRedirect, res.StatusCode)
					return nil
				}
			})
			grp.Go(func() error {
				internal.RunContext(ctx, t, tt, nil, command(cfg))
				return nil
			})
			a.NoError(grp.Wait())
		})
	}
}
