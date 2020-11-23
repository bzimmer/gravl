package gravl

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"

	"github.com/bzimmer/gravl/pkg/common"
	"github.com/bzimmer/gravl/pkg/strava"
	"github.com/bzimmer/gravl/pkg/wta"
)

func newRouter(client *wta.Client) *gin.Engine {
	r := gin.New()
	r.Use(
		gin.Recovery(),
		common.LogMiddleware(),
		sessions.Sessions("session", gothic.Store.(sessions.Store)))

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl.html", nil)
	})
	r.LoadHTMLGlob("web/views/templates/*.tmpl.html")
	r.Static("/static", "web/static")

	r.GET("/version/", common.VersionHandler())

	p := r.Group("/wta")
	p.GET("/regions/", wta.RegionsHandler(client))
	p.GET("/reports/", wta.TripReportsHandler(client))
	p.GET("/reports/:reporter", wta.TripReportsHandler(client))

	p = r.Group("/strava")
	p.GET("/auth/login", strava.AuthHandler())
	p.GET("/auth/callback", strava.AuthCallbackHandler())

	return r
}

var serveCommand = &cli.Command{
	Name:     "serve",
	Category: "api",
	Usage:    "REST endpoints",
	Flags: []cli.Flag{
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:  "strava.api-key",
			Usage: "API key for Strava API",
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:  "strava.api-secret",
			Usage: "API secret for Strava API",
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:  "serve.origin",
			Value: "http://localhost",
			Usage: "Callback origin",
		}),
		altsrc.NewIntFlag(&cli.IntFlag{
			Name:  "serve.port",
			Value: 8080,
			Usage: "Port on which to listen",
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:  "serve.session-key",
			Usage: "Session key",
		}),
	},
	Before: func(c *cli.Context) error {
		callback := fmt.Sprintf("%s:%d/strava/auth/callback", c.String("serve.origin"), c.Int("serve.port"))
		log.Info().Str("callback", callback).Msg("preparing to serve")
		goth.UseProviders(newStravaAuthProvider(c, callback))
		gothic.Store = cookie.NewStore([]byte(c.String("serve.session-key")))
		gothic.GetProviderName = func(req *http.Request) (string, error) {
			return "strava", nil
		}
		return nil
	},
	Action: func(c *cli.Context) error {
		client, err := wta.NewClient(wta.WithHTTPTracing(c.Bool("http-tracing")))
		if err != nil {
			return err
		}
		r := newRouter(client)
		log.Info().Msg("serving ...")
		address := fmt.Sprintf("0.0.0.0:%d", c.Int("serve.port"))
		if err := r.Run(address); err != nil {
			return err
		}
		return nil
	},
}
