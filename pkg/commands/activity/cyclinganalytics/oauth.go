package cyclinganalytics

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"golang.org/x/oauth2"

	ca "github.com/bzimmer/gravl/pkg/activity/cyclinganalytics"
	"github.com/bzimmer/gravl/pkg/commands"
	"github.com/bzimmer/gravl/pkg/web"
)

func newRouter(c *cli.Context) (*gin.Engine, error) {
	r := gin.New()
	r.Use(gin.Recovery(), web.LogMiddleware())
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusTemporaryRedirect, "/cyclinganalytics/auth/login")
	})
	state, err := commands.Token(16)
	if err != nil {
		return nil, err
	}
	address := fmt.Sprintf("%s:%d", c.String("origin"), c.Int("port"))
	p := r.Group("/cyclinganalytics")
	config := &oauth2.Config{
		ClientID:     c.String("cyclinganalytics.client-id"),
		ClientSecret: c.String("cyclinganalytics.client-secret"),
		Scopes:       []string{"read_account,read_email,read_athlete,read_rides,create_rides"},
		RedirectURL:  fmt.Sprintf("%s/cyclinganalytics/auth/callback", address),
		Endpoint:     ca.Endpoint}
	p.GET("/auth/login", gin.WrapF(web.AuthHandler(config, state)))
	p.GET("/auth/callback", gin.WrapF(web.AuthCallbackHandler(config, state)))
	return r, nil
}

var oauthCommand = &cli.Command{
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
		router, err := newRouter(c)
		if err != nil {
			return err
		}
		address := fmt.Sprintf("0.0.0.0:%d", c.Int("port"))
		log.Info().Str("address", address).Msg("serving ...")
		return router.Run(address)
	},
}