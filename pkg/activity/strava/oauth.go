package strava

import (
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"golang.org/x/oauth2"

	st "github.com/bzimmer/activity/strava"
	"github.com/bzimmer/gravl/pkg"
	"github.com/bzimmer/gravl/pkg/web"
)

func newRouter(c *cli.Context) (*http.ServeMux, error) {
	state, err := pkg.Token(16)
	if err != nil {
		return nil, err
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		http.Redirect(w, r, "/strava/auth/login", http.StatusTemporaryRedirect)
	})
	address := fmt.Sprintf("%s:%d", c.String("origin"), c.Int("port"))
	config := &oauth2.Config{
		ClientID:     c.String("strava.client-id"),
		ClientSecret: c.String("strava.client-secret"),
		Scopes:       []string{"read_all,profile:read_all,activity:read_all,activity:write"},
		RedirectURL:  fmt.Sprintf("%s/strava/auth/callback", address),
		Endpoint:     st.Endpoint()}
	handle := web.NewLogHandler(&log.Logger)
	mux.Handle("/strava/auth/login", handle(web.AuthHandler(config, state)))
	mux.Handle("/strava/auth/callback", handle(web.AuthCallbackHandler(config, state)))
	return mux, nil
}

func oauthCommand() *cli.Command {
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
				Value: 9001,
				Usage: "Port on which to listen",
			},
		},
		Action: func(c *cli.Context) error {
			mux, err := newRouter(c)
			if err != nil {
				return err
			}
			address := fmt.Sprintf("0.0.0.0:%d", c.Int("port"))
			log.Info().Str("address", address).Msg("serving")
			return http.ListenAndServe(address, mux)
		},
	}
}
