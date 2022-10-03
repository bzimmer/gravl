package activity

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"golang.org/x/oauth2"
	"golang.org/x/sync/errgroup"

	"github.com/bzimmer/gravl"
	"github.com/bzimmer/gravl/web"
)

// OAuthConfig provides configuration options
type OAuthConfig struct {
	// Port on which to listen
	Port int
	// Provider at which to authenticate
	Provider string
	// Scopes to request with the credentials
	Scopes []string
	// RedirectURL on successful authentication
	RedirectURL string
	// Started is the channel on which the server url is communicated
	Started chan<- *url.URL
}

func newHandler(c *cli.Context, cfg *OAuthConfig) (http.Handler, error) {
	state, err := gravl.Token(16)
	if err != nil {
		return nil, err
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		http.Redirect(w, r, fmt.Sprintf("/%s/auth/login", cfg.Provider), http.StatusTemporaryRedirect)
	})
	address := fmt.Sprintf("%s:%d", c.String("origin"), c.Int("port"))
	redirectURL := cfg.RedirectURL
	if redirectURL == "" {
		redirectURL = fmt.Sprintf("%s/%s/auth/callback", address, cfg.Provider)
	}
	config := &oauth2.Config{
		ClientID:     c.String(fmt.Sprintf("%s-client-id", cfg.Provider)),
		ClientSecret: c.String(fmt.Sprintf("%s-client-secret", cfg.Provider)),
		Scopes:       cfg.Scopes,
		RedirectURL:  redirectURL,
		Endpoint:     gravl.Runtime(c).Endpoints[cfg.Provider]}
	log.Info().Str("redirect", config.RedirectURL).Msg(c.Command.Name)
	handle := web.NewLogHandler(&log.Logger)
	mux.Handle(fmt.Sprintf("/%s/auth/login", cfg.Provider), handle(web.AuthHandler(config, state)))
	mux.Handle(fmt.Sprintf("/%s/auth/callback", cfg.Provider), handle(web.AuthCallbackHandler(config, state)))
	return mux, nil
}

func newListener(cfg *OAuthConfig) (net.Listener, error) {
	listener, err := net.Listen("tcp4", fmt.Sprintf("127.0.0.1:%d", cfg.Port))
	if err != nil {
		return nil, err
	}
	return listener, nil
}

func oauth(c *cli.Context, cfg *OAuthConfig) error {
	mux, err := newHandler(c, cfg)
	if err != nil {
		return err
	}
	listener, err := newListener(cfg)
	if err != nil {
		return err
	}
	svr := &http.Server{
		Handler:           mux,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
	}
	grp, ctx := errgroup.WithContext(c.Context)
	grp.Go(func() error {
		if err = svr.Serve(listener); !errors.Is(err, http.ErrServerClosed) {
			log.Info().Err(err).Msg("closed")
			return err
		}
		return nil
	})
	grp.Go(func() error {
		<-ctx.Done()
		return svr.Close()
	})
	grp.Go(func() error {
		if cfg.Started == nil {
			return nil
		}
		var u *url.URL
		u, err = url.Parse("http://" + listener.Addr().String())
		if err != nil {
			return err
		}
		log.Info().Str("address", u.String()).Msg("serving")
		select {
		case <-c.Context.Done():
			return c.Context.Err()
		case cfg.Started <- u:
			return nil
		}
	})
	if err = grp.Wait(); err != nil {
		if errors.Is(err, context.Canceled) {
			return nil
		}
		return err
	}
	return nil
}

func OAuthCommand(cfg *OAuthConfig) *cli.Command {
	return &cli.Command{
		Name:  "oauth",
		Usage: "Authentication endpoints for access and refresh tokens",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "origin",
				Value: "http://localhost",
				Usage: "Callback origin",
			},
			&cli.IntFlag{
				Name:  "port",
				Value: cfg.Port,
				Usage: "Port on which to listen",
			},
		},
		Action: func(c *cli.Context) error {
			return oauth(c, cfg)
		},
	}
}
