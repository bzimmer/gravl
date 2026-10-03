package activity_test

import (
	"context"
	"net"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"
	"golang.org/x/sync/errgroup"

	"github.com/bzimmer/gravl/activity"
	"github.com/bzimmer/gravl/internal"
)

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
			ctx, cancel := context.WithTimeout(t.Context(), time.Second)
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
					req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), http.NoBody)
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
				internal.RunContext(ctx, t, tt, nil, func(_ *testing.T, _ string) *cli.Command {
					return activity.OAuthCommand(cfg)
				})
				return nil
			})
			a.NoError(grp.Wait())
		})
	}
}

func runOAuth(ctx context.Context, t *testing.T, tt *internal.Harness, cfg *activity.OAuthConfig) {
	internal.RunContext(ctx, t, tt, nil, func(_ *testing.T, _ string) *cli.Command {
		return activity.OAuthCommand(cfg)
	})
}

func TestOAuthPortInUse(t *testing.T) {
	a := assert.New(t)
	var lc net.ListenConfig
	busy, err := lc.Listen(t.Context(), "tcp4", "127.0.0.1:0")
	a.NoError(err)
	defer busy.Close()
	port := busy.Addr().(*net.TCPAddr).Port //nolint:errcheck // a tcp4 listener's address
	tt := &internal.Harness{Name: "port in use", Args: []string{"test", "oauth"}, Err: "address already in use"}
	runOAuth(t.Context(), t, tt, &activity.OAuthConfig{Provider: "foobar", Port: port})
}

func TestOAuthShutdown(t *testing.T) {
	tests := []struct {
		name    string
		started bool
		ctx     func(context.Context) (context.Context, context.CancelFunc)
		err     string
	}{
		{
			// nothing reports the address, so the server runs until the context ends
			name: "no started channel",
			ctx: func(ctx context.Context) (context.Context, context.CancelFunc) {
				return context.WithTimeout(ctx, 100*time.Millisecond)
			},
		},
		{
			// cancelled while waiting to report the address: a clean shutdown
			name:    "cancelled before the address is read",
			started: true,
			ctx: func(ctx context.Context) (context.Context, context.CancelFunc) {
				ctx, cancel := context.WithCancel(ctx)
				time.AfterFunc(100*time.Millisecond, cancel)
				return ctx, cancel
			},
		},
		{
			// a deadline is not a cancellation, so it is reported
			name:    "deadline before the address is read",
			started: true,
			ctx: func(ctx context.Context) (context.Context, context.CancelFunc) {
				return context.WithTimeout(ctx, 100*time.Millisecond)
			},
			err: "deadline exceeded",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := tc.ctx(t.Context())
			defer cancel()
			cfg := &activity.OAuthConfig{Provider: "foobar"}
			if tc.started {
				// unbuffered and never read, so reporting the address blocks
				cfg.Started = make(chan *url.URL)
			}
			tt := &internal.Harness{Name: tc.name, Args: []string{"test", "oauth"}, Err: tc.err}
			runOAuth(ctx, t, tt, cfg)
		})
	}
}
