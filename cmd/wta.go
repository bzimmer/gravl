package cmd

import (
	"github.com/bzimmer/wta/pkg/common"
	w "github.com/bzimmer/wta/pkg/wta"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func wta(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		// query the most recent if no reporter specified
		args = append(args, "")
	}

	c, err := w.NewClient()
	if err != nil {
		return err
	}
	reports := make([]*w.TripReport, 0)
	for _, arg := range args {
		tr, err := c.Reports.TripReports(cmd.Context(), arg)
		if err != nil {
			return err
		}
		for _, r := range tr {
			reports = append(reports, r)
		}
	}

	encoder := common.NewEncoder(compact)
	err = encoder.Encode(reports)
	if err != nil {
		return err
	}
	return nil
}

func init() {
	rootCmd.AddCommand(wtaCmd)
}

var wtaCmd = &cobra.Command{
	Use:     "wta",
	Short:   "Run wta",
	Long:    `Run wta`,
	Aliases: []string{"w"},
	PreRun: func(cmd *cobra.Command, args []string) {
		log.Info().
			Str("url", "https://www.wta.org/").
			Msg("Please support the WTA")
	},
	RunE: wta,
}