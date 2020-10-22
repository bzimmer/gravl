package cmd

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/bzimmer/wta/pkg/common"
	na "github.com/bzimmer/wta/pkg/noaa"

	gj "github.com/paulmach/go.geojson"
	"github.com/spf13/cobra"
)

func noaa(cmd *cobra.Command, args []string) error {
	c, err := na.NewClient(
		na.WithTimeout(10 * time.Second),
	)
	if err != nil {
		return err
	}

	var forecast *na.Forecast
	switch len(args) {
	case 0:
		geom := &gj.Geometry{}
		decoder := common.NewDecoder()
		err := decoder.Decode(geom)
		if err != nil {
			return err
		}
		if !geom.IsPoint() {
			return errors.New("only Point geometries are supported")
		}
		forecast, err = c.Points.Forecast(cmd.Context(), geom.Point[1], geom.Point[0])
	case 2:
		lat, err := strconv.ParseFloat(args[0], 64)
		if err != nil {
			return err
		}
		lng, err := strconv.ParseFloat(args[1], 64)
		if err != nil {
			return err
		}
		forecast, err = c.Points.Forecast(cmd.Context(), lat, lng)
	case 3:
		wfo := args[0]
		x, err := strconv.Atoi(args[1])
		if err != nil {
			return err
		}
		y, err := strconv.Atoi(args[2])
		if err != nil {
			return err
		}
		forecast, err = c.GridPoints.Forecast(cmd.Context(), wfo, x, y)
	default:
		return fmt.Errorf("only 2 or 3 arguments allowed [%v]", args)
	}
	if err != nil {
		return err
	}
	encoder := common.NewEncoder(compact)
	err = encoder.Encode(forecast)
	if err != nil {
		return err
	}
	return nil
}

func init() {
	rootCmd.AddCommand(noaaCmd)
}

var noaaCmd = &cobra.Command{
	Use:     "noaa",
	Short:   "Run noaa",
	Long:    `Run noaa`,
	Aliases: []string{"n", "f"},
	RunE:    noaa,
}
