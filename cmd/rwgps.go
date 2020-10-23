package cmd

import (
	"errors"
	"strconv"

	"github.com/bzimmer/wta/pkg/common"
	rw "github.com/bzimmer/wta/pkg/rwgps"

	gj "github.com/paulmach/go.geojson"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

var (
	trip      bool
	route     bool
	apiKey    string
	authToken string
)

func rwgps(cmd *cobra.Command, args []string) error {
	var (
		err error
		tr  *gj.FeatureCollection
	)
	// nothing to query
	if len(args) == 0 {
		return nil
	}
	lvl, err := zerolog.ParseLevel(verbosity)
	if err != nil {
		return err
	}
	c, err := rw.NewClient(
		rw.WithAPIKey(apiKey),
		rw.WithAuthToken(authToken),
		rw.WithVerboseLogging(lvl == zerolog.DebugLevel),
	)
	if err != nil {
		return err
	}
	encoder := common.NewEncoder(compact)
	for _, arg := range args {
		x, err := strconv.ParseInt(arg, 0, 0)
		if err != nil {
			return err
		}
		if trip {
			tr, err = c.Trips.Trip(cmd.Context(), x)
		} else {
			tr, err = c.Trips.Route(cmd.Context(), x)
		}
		if err != nil {
			return err
		}
		err = encoder.Encode(tr)
		if err != nil {
			return err
		}
	}
	return nil
}

func init() {
	rootCmd.AddCommand(rwgpsCmd)
	rwgpsCmd.Flags().StringVarP(&apiKey, "rwgps_apikey", "k", "", "API key")
	rwgpsCmd.Flags().StringVarP(&authToken, "rwgps_authtoken", "a", "", "Auth token")

	rwgpsCmd.Flags().BoolVarP(&trip, "trip", "t", false, "Trip")
	rwgpsCmd.Flags().BoolVarP(&route, "route", "r", false, "Route")
}

var rwgpsCmd = &cobra.Command{
	Use:     "rwgps",
	Short:   "Run rwgps",
	Long:    `Run rwgps`,
	Aliases: []string{"r"},
	RunE:    rwgps,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if !(trip || route) {
			return errors.New("one of trip or route must chosen")
		}
		if trip && route {
			return errors.New("only one of trip or route allowed")
		}
		return nil
	},
}
