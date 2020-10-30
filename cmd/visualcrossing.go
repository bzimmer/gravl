package cmd

import (
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"

	"github.com/bzimmer/gravl/pkg/common"
	"github.com/bzimmer/gravl/pkg/common/wx"
	vc "github.com/bzimmer/gravl/pkg/visualcrossing"
)

var (
	visualcrossingAPIKey string
)

func visualcrossing(cmd *cobra.Command, args []string) error {
	var (
		err  error
		fcst []*wx.Forecast
	)
	lvl, err := zerolog.ParseLevel(verbosity)
	if err != nil {
		return err
	}

	c, err := vc.NewClient(
		vc.WithAPIKey(visualcrossingAPIKey),
		vc.WithVerboseLogging(lvl == zerolog.DebugLevel),
	)
	if err != nil {
		return err
	}

	fcst, err = c.Forecast.Forecast(cmd.Context(),
		vc.WithLocations(args...),
		vc.WithAstronomy(true),
		vc.WithUnits(vc.UnitsUS),
		vc.WithAggregateHours(12),
		vc.WithAlerts(vc.AlertLevelDetail))
	if err != nil {
		return err
	}

	fc, err := wx.NewFeatureCollection(fcst...)
	if err != nil {
		return err
	}

	encoder := common.NewEncoder(compact)
	err = encoder.Encode(fc)
	if err != nil {
		return err
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
