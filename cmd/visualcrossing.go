package cmd

import (
	"github.com/spf13/cobra"

	"github.com/bzimmer/gravl/pkg/common/wx"
	vc "github.com/bzimmer/gravl/pkg/visualcrossing"
)

var (
	visualcrossingAPIKey string
)

func visualcrossing(cmd *cobra.Command, args []string) error {
	c, err := vc.NewClient(
		vc.WithAPIKey(visualcrossingAPIKey),
		vc.WithVerboseLogging(debug),
	)
	if err != nil {
		return err
	}

	for _, arg := range args {
		fcst, err := c.Forecast.Forecast(
			cmd.Context(),
			vc.WithLocation(arg),
			vc.WithAstronomy(true),
			vc.WithUnits(vc.UnitsUS),
			vc.WithAggregateHours(12),
			vc.WithAlerts(vc.AlertLevelDetail))
		if err != nil {
			return err
		}
		fc, err := wx.NewFeatureCollection(fcst)
		if err != nil {
			return err
		}
		err = encoder.Encode(fc)
		if err != nil {
			return err
		}
	}
	return nil
}

func init() {
	rootCmd.AddCommand(visualcrossingCmd)
	visualcrossingCmd.Flags().StringVarP(&visualcrossingAPIKey, "visualcrossing_apikey", "k", "", "API key")
}

var visualcrossingCmd = &cobra.Command{
	Use:     "vc",
	Short:   "Run visualcrossing",
	Long:    `Run visualcrossing`,
	Aliases: []string{""},
	RunE:    visualcrossing,
}
