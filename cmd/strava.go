package cmd

import (
	"github.com/markbates/goth"
	au "github.com/markbates/goth/providers/strava"
	"github.com/spf13/cobra"

	sa "github.com/bzimmer/gravl/pkg/strava"
)

func newStravaAuthProvider(callback string) goth.Provider {
	return au.New(
		stravaAPIKey, stravaAPISecret, callback,
		// appears to be a bug where scope varargs do not work properly
		"read_all,profile:read_all,activity:read_all")
}

func strava(cmd *cobra.Command, args []string) error {
	c, err := sa.NewClient(
		sa.WithVerboseLogging(debug),
		sa.WithAPICredentials(stravaAccessToken, stravaRefreshToken),
		sa.WithProvider(newStravaAuthProvider("")))
	if err != nil {
		return err
	}
	if athlete {
		ath, err := c.Athlete.Athlete(cmd.Context())
		if err != nil {
			return err
		}
		encoder.Encode(ath)
	}
	if refresh {
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

	stravaCmd.Flags().BoolVarP(&athlete, "athlete", "a", false, "display the authenticated athlete")
	stravaCmd.Flags().BoolVarP(&refresh, "refresh", "r", false, "refresh the access token")
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
