package cyclinganalytics

import (
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"golang.org/x/oauth2"

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
		http.Redirect(w, r, "/cyclinganalytics/auth/login", http.StatusTemporaryRedirect)
	})
	address := fmt.Sprintf("%s:%d", c.String("origin"), c.Int("port"))
	config := &oauth2.Config{
		ClientID:     c.String("cyclinganalytics-client-id"),
		ClientSecret: c.String("cyclinganalytics-client-secret"),
		Scopes:       []string{"read_account,read_email,read_athlete,read_rides,create_rides"},
		RedirectURL:  fmt.Sprintf("%s/cyclinganalytics/auth/callback", address),
		Endpoint:     pkg.Runtime(c).Endpoints[Provider]}
	// The redirect url is not dynamic, it must be configured on the cyclinganalytics website
	//  https://www.cyclinganalytics.com/account/apps
	// in order for the flow to work correctly. If the redirect url at ca is not exactly the
	// same (hostname, port, path) the redirect will fail.
	log.Info().Str("redirect", config.RedirectURL).Msg(c.Command.Name)
	handle := web.NewLogHandler(&log.Logger)
	mux.Handle("/cyclinganalytics/auth/login", handle(web.AuthHandler(config, state)))
	mux.Handle("/cyclinganalytics/auth/callback", handle(web.AuthCallbackHandler(config, state)))
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
				Value: 9002,
				Usage: "Port on which to listen",
			},
		},
		Action: func(c *cli.Context) error {
			mux, err := newRouter(c)
			if err != nil {
				return err
			}
			address := fmt.Sprintf("%s:%d", c.String("origin"), c.Int("port"))
			log.Info().Str("address", address).Msg("serving")
			return http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", c.Int("port")), mux)
		},
	}
}
