package cmd

import (
	"github.com/spf13/cobra"

	sa "github.com/bzimmer/gravl/pkg/strava"
)

var (
	stravaAPIKey       string
	stravaAPISecret    string
	stravaAccessToken  string
	stravaRefreshToken string
)

func strava(cmd *cobra.Command, args []string) error {
	c, err := sa.NewClient(
		sa.WithVerboseLogging(debug),
		sa.WithAPICredentials(stravaAccessToken, stravaRefreshToken))
	if err != nil {
		return err
	}
	ath, err := c.Athlete.Athlete(cmd.Context())
	if err != nil {
		return err
	}
	encoder.Encode(ath)
	return nil
}

func init() {
	rootCmd.AddCommand(stravaCmd)
	stravaCmd.Flags().StringVarP(&stravaAPIKey, "strava_key", "", "", "API key")
	stravaCmd.Flags().StringVarP(&stravaAPISecret, "strava_secret", "", "", "API secret")
	stravaCmd.Flags().StringVarP(&stravaAccessToken, "strava_access_token", "", "", "API access token")
	stravaCmd.Flags().StringVarP(&stravaRefreshToken, "strava_refresh_token", "", "", "API refresh token")
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
