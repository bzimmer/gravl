package cmd

import (
	"net/http"
	"strconv"
	"time"

	"github.com/markbates/goth"
	au "github.com/markbates/goth/providers/strava"
	"github.com/spf13/cobra"

	"github.com/bzimmer/gravl/pkg/common"
	sa "github.com/bzimmer/gravl/pkg/strava"
)

var (
	streams = []string{
		"latlng", "altitude", "distance",
	}
)

func newStravaAuthProvider(callback string) goth.Provider {
	provider := au.New(
		stravaAPIKey, stravaAPISecret, callback,
		// appears to be a bug where scope varargs do not work properly
		"read_all,profile:read_all,activity:read_all")
	transport := http.DefaultTransport
	if httptracing {
		transport = &common.VerboseTransport{
			Transport: transport,
		}
	}
	provider.HTTPClient = &http.Client{
		Timeout:   10 * time.Second,
		Transport: transport,
	}
	return provider
}

func strava(cmd *cobra.Command, args []string) error {
	c, err := sa.NewClient(
		sa.WithHTTPTracing(httptracing),
		sa.WithAPICredentials(stravaAccessToken, stravaRefreshToken),
		sa.WithProvider(newStravaAuthProvider("")))
	if err != nil {
		return err
	}
	if stravaActivity {
		for _, arg := range args {
			activityID, err := strconv.ParseInt(arg, 0, 64)
			streams, err := c.Activity.Streams(cmd.Context(), activityID, streams...)
			if err != nil {
				return err
			}
			encoder.Encode(streams)
		}
	}
	if stravaAthlete {
		ath, err := c.Athlete.Athlete(cmd.Context())
		if err != nil {
			return err
		}
		encoder.Encode(ath)
	}
	if stravaRefresh {
		tokens, err := c.Auth.Refresh(cmd.Context())
		if err != nil {
			return err
		}
		encoder.Encode(tokens)
	}
	return nil
}

func init() {
	rootCmd.AddCommand(stravaCmd)
	stravaCmd.Flags().StringVarP(&stravaAPIKey, "strava_key", "", "", "API key")
	stravaCmd.Flags().StringVarP(&stravaAPISecret, "strava_secret", "", "", "API secret")
	stravaCmd.Flags().StringVarP(&stravaAccessToken, "strava_access_token", "", "", "API access token")
	stravaCmd.Flags().StringVarP(&stravaRefreshToken, "strava_refresh_token", "", "", "API refresh token")

	stravaCmd.Flags().BoolVarP(&stravaActivity, "activity", "a", false, "Return stream data for the activity")
	stravaCmd.Flags().BoolVarP(&stravaAthlete, "athlete", "u", false, "Display the authenticated athlete")
	stravaCmd.Flags().BoolVarP(&stravaRefresh, "refresh", "r", false, "Refresh the access token")
}

var stravaCmd = &cobra.Command{
	Use:     "strava",
	Short:   "Run strava",
	Long:    `Run strava`,
	Aliases: []string{"a"},
	RunE:    strava,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}
