package cmd

import (
	"fmt"

	w "github.com/bzimmer/wta/pkg/wta"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	port int
)

func serve(cmd *cobra.Command, args []string) error {
	log.Info().Msg("configuring to serve")
	c, err := w.NewClient()
	if err != nil {
		return err
	}
	r := w.NewRouter(c)

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
