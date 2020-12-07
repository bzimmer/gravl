package gravl

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
	"golang.org/x/oauth2"

	"github.com/bzimmer/gravl/pkg/common/web"
	"github.com/bzimmer/gravl/pkg/cyclinganalytics"
	"github.com/bzimmer/gravl/pkg/strava"
)

var index = []byte(`
<html>
	<head><title>Gravl</title></head>
	<body>
		<ul>
		<li><a href="/strava/auth/login">Auth with Strava</a></li>
		<li><a href="/cyclinganalytics/auth/login">Auth with Cycling Analytics</a></li>
		</ul>
	</body>
</html>`)

func newRouter(c *cli.Context) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery(), web.LogMiddleware())

	r.GET("/", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html", index)
	})

	r.GET("/version/", gin.WrapF(web.VersionHandler()))

	state := mustRandomString(16)
	address := fmt.Sprintf("%s:%d", c.String("serve.origin"), c.Int("serve.port"))

	p := r.Group("/cyclinganalytics")
	config := &oauth2.Config{
		ClientID:     c.String("cyclinganalytics.client-id"),
		ClientSecret: c.String("cyclinganalytics.client-secret"),
		Scopes:       []string{"read_account,read_email,read_athlete,read_rides"},
		RedirectURL:  fmt.Sprintf("%s/cyclinganalytics/auth/callback", address),
		Endpoint:     cyclinganalytics.Endpoint}
	p.GET("/auth/login", gin.WrapF(web.AuthHandler(config, state)))
	p.GET("/auth/callback", gin.WrapF(web.AuthCallbackHandler(config, state)))

	p = r.Group("/strava")
	config = &oauth2.Config{
		ClientID:     c.String("strava.client-id"),
		ClientSecret: c.String("strava.client-secret"),
		Scopes:       []string{"read_all,profile:read_all,activity:read_all"},
		RedirectURL:  fmt.Sprintf("%s/strava/auth/callback", address),
		Endpoint:     strava.Endpoint}
	p.GET("/auth/login", gin.WrapF(web.AuthHandler(config, state)))
	p.GET("/auth/callback", gin.WrapF(web.AuthCallbackHandler(config, state)))

	return r
}

var serveCommand = &cli.Command{
	Name:     "serve",
	Category: "api",
	Usage:    "REST endpoints",
	Flags: merge([]cli.Flag{
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:  "serve.origin",
			Value: "http://localhost",
			Usage: "Callback origin",
		}),
		altsrc.NewIntFlag(&cli.IntFlag{
			Name:  "serve.port",
			Value: 8080,
			Usage: "Port on which to listen",
		})},
		cyclingAnalyticsAuthFlags,
		stravaAuthFlags,
	),
	Action: func(c *cli.Context) error {
		r := newRouter(c)
		address := fmt.Sprintf("0.0.0.0:%d", c.Int("serve.port"))
		log.Info().Str("address", address).Msg("serving ...")
		if err := r.Run(address); err != nil {
			return err
		}
		return nil
	},
}
