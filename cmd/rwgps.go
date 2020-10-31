package cmd

import (
	"errors"
	"strconv"

	gj "github.com/paulmach/go.geojson"
	"github.com/spf13/cobra"

	rw "github.com/bzimmer/gravl/pkg/rwgps"
)

func rwgps(cmd *cobra.Command, args []string) error {
	var (
		err error
		tr  *gj.FeatureCollection
	)
	c, err := rw.NewClient(
		rw.WithAPIKey(rwgpsAPIKey),
		rw.WithAuthToken(rwgpsAuthToken),
		rw.WithVerboseLogging(debug),
	)
	if err != nil {
		return err
	}
	if rwgpsUser {
		user, err := c.Users.AuthenticatedUser(cmd.Context())
		if err != nil {
			return err
		}
		err = encoder.Encode(user)
		if err != nil {
			return err
		}
		return nil
	}
	for _, arg := range args {
		x, err := strconv.ParseInt(arg, 0, 0)
		if err != nil {
			return err
		}
		if rwgpsTrip {
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
	rwgpsCmd.Flags().StringVarP(&rwgpsAPIKey, "rwgps_apikey", "k", "", "API key")
	rwgpsCmd.Flags().StringVarP(&rwgpsAuthToken, "rwgps_authtoken", "a", "", "Auth token")

	rwgpsCmd.Flags().BoolVarP(&rwgpsTrip, "trip", "t", false, "Trip")
	rwgpsCmd.Flags().BoolVarP(&rwgpsRoute, "route", "r", false, "Route")
	rwgpsCmd.Flags().BoolVarP(&rwgpsUser, "user", "u", false, "User")
}

var rwgpsCmd = &cobra.Command{
	Use:     "rwgps",
	Short:   "Run rwgps",
	Long:    `Run rwgps`,
	Aliases: []string{"r"},
	RunE:    rwgps,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if rwgpsTrip && rwgpsRoute {
			return errors.New("only one of trip or route allowed")
		}
		return nil
	},
}
