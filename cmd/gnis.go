package cmd

import (
	"github.com/bzimmer/wta/pkg/common"
	gn "github.com/bzimmer/wta/pkg/gnis"

	"github.com/spf13/cobra"
)

func gnis(cmd *cobra.Command, args []string) error {
	g, err := gn.NewClient()
	if err != nil {
		return err
	}
	encoder := common.NewEncoder(compact)
	for _, arg := range args {
		features, err := g.GeoNames.Query(cmd.Context(), arg)
		if err != nil {
			return err
		}
		err = encoder.Encode(features)
		if err != nil {
			return err
		}
	}
	return nil
}

func init() {
	rootCmd.AddCommand(gnisCmd)
}

var gnisCmd = &cobra.Command{
	Use:     "gnis",
	Short:   "Run gnis",
	Long:    `Run gnis`,
	Aliases: []string{"g"},
	RunE:    gnis,
}
