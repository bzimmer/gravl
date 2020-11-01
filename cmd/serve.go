package cmd

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/bzimmer/gravl/pkg/common"
	sa "github.com/bzimmer/gravl/pkg/strava"
	wa "github.com/bzimmer/gravl/pkg/wta"
)

func newRouter(client *wa.Client) *gin.Engine {
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
	p.GET("/regions/", wa.RegionsHandler())
	p.GET("/reports/", wa.TripReportsHandler(client))
	p.GET("/reports/:reporter", wa.TripReportsHandler(client))

	p = r.Group("/strava")
	p.GET("/auth/login", sa.AuthHandler())
	p.GET("/auth/callback", sa.AuthCallbackHandler())

	return r
}

func serve(cmd *cobra.Command, args []string) error {
	log.Info().Msg("configuring to serve")
	c, err := wa.NewClient()
	if err != nil {
		return err
	}
	r := newRouter(c)
	address := fmt.Sprintf("0.0.0.0:%d", port)
	if err := r.Run(address); err != nil {
		return err
	}
	return nil
}

func initAuth(cmd *cobra.Command) error {
	callback := fmt.Sprintf("%s:%d/strava/auth/callback", origin, port)
	log.Info().Str("callback", callback).Msg("initAuth")
	goth.UseProviders(newStravaAuthProvider(callback))
	gothic.Store = cookie.NewStore([]byte(sessionKey))
	gothic.GetProviderName = func(req *http.Request) (string, error) {
		return "strava", nil
	}
	return nil
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().IntVarP(&port, "port", "p", 1122, "Port on which to listen")
	serveCmd.Flags().StringVarP(&origin, "origin", "", "", "Origin URL")
	serveCmd.Flags().StringVarP(&sessionKey, "session_key", "", "", "session key")
	serveCmd.Flags().StringVarP(&stravaAPIKey, "strava_key", "", "", "API key")
	serveCmd.Flags().StringVarP(&stravaAPISecret, "strava_secret", "", "", "API secret")
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run the server",
	Long:  `Run the server`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := initAuth(cmd); err != nil {
			return nil
		}
		if err := serve(cmd, args); err != nil {
			return err
		}
		return nil
	},
}
