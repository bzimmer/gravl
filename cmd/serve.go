package cmd

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/bzimmer/gravl/pkg/common"
	wa "github.com/bzimmer/gravl/pkg/wta"
)

var (
	port int
)

func newRouter(client *wa.Client) *gin.Engine {
	r := gin.New()
	r.Use(common.LogMiddleware(), gin.Recovery())
	r.GET("/version/", common.VersionHandler())
	r.GET("/regions/", wa.RegionsHandler())
	r.GET("/reports/", wa.TripReportsHandler(client))
	r.GET("/reports/:reporter", wa.TripReportsHandler(client))
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

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().IntVarP(&port, "port", "p", 1122, "Port on which to listen")
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run the server",
	Long:  `Run the server`,
	RunE:  serve,
}
