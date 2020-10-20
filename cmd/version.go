package cmd

import (
	"fmt"

	"github.com/bzimmer/wta/pkg"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:     "version",
	Short:   "Print the version number of Hugo",
	Long:    `All software has versions. This is Hugo's`,
	Aliases: []string{"v"},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(`{"version": "` + pkg.BuildVersion + `"}`)
	},
}