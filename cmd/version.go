package cmd

import (
	"github.com/spf13/cobra"

	"github.com/bzimmer/gravl/pkg"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:     "version",
	Short:   "Print the version number of Hugo",
	Long:    `All software has versions. This is Hugo's`,
	Aliases: []string{"v"},
	RunE: func(cmd *cobra.Command, args []string) error {
		err := encoder.Encode(map[string]string{"version": pkg.BuildVersion})
		if err != nil {
			return err
		}
		return nil
	},
}
