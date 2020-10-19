package cmd

import (
	"fmt"
	"strconv"
	"time"

	"github.com/bzimmer/wta/pkg/common"
	na "github.com/bzimmer/wta/pkg/noaa"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func noaa(cmd *cobra.Command, args []string) error {
	c, err := na.NewClient(
		na.WithTimeout(10 * time.Second),
	)
	if err != nil {
		return err
	}

	var (
		// point    *na.GridPoint
		forecast *na.Forecast
	)
	switch len(args) {
	case 2:
		lat, err := strconv.ParseFloat(args[0], 64)
		if err != nil {
			log.Error().Err(err).Send()
		}
		lng, err := strconv.ParseFloat(args[1], 64)
		if err != nil {
			log.Error().Err(err).Send()
		}
		// if ctx.IsSet("point") {
		// 	point, err = c.Points.GridPoint(cmd.Context(), lat, lng)
		// }
		forecast, err = c.Points.Forecast(cmd.Context(), lat, lng)
	case 3:
		// check for -p and err if true
		wfo := args[0]
		x, _ := strconv.Atoi(args[1])
		y, _ := strconv.Atoi(args[2])
		forecast, err = c.GridPoints.Forecast(cmd.Context(), wfo, x, y)
	default:
		return fmt.Errorf("only 2 or 3 arguments allowed")
	}
	if err != nil {
		return err
	}
	encoder := common.NewEncoder(compact)
	// if ctx.IsSet("point") {
	// 	err = encoder.Encode(point)
	// } else {
	err = encoder.Encode(forecast)
	// }
	if err != nil {
		return err
	}
	return nil
}

func init() {
	rootCmd.AddCommand(noaaCmd)
}

var noaaCmd = &cobra.Command{
	Use:   "noaa",
	Short: "Run noaa",
	Long:  `Run noaa`,
	RunE:  noaa,
}
