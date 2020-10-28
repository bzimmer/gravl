package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"github.com/spf13/cobra"
	"github.com/tormoder/fit"
)

func runfit(cmd *cobra.Command, args []string) error {
	// Read our FIT test file data
	// testFile := filepath.Join("testdata", "fitsdk", "Activity.fit")
	testFile := args[0]
	testData, err := ioutil.ReadFile(testFile)
	if err != nil {
		return err
	}

	// Decode the FIT file data
	fit, err := fit.Decode(bytes.NewReader(testData))
	if err != nil {
		return err
	}

	// Inspect the TimeCreated field in the FileId message
	fmt.Println(fit.FileId.TimeCreated)

	// Inspect the dynamic Product field in the FileId message
	fmt.Println(fit.FileId.GetProduct())

	// Inspect the FIT file type
	fmt.Println(fit.Type())

	// Get the actual activity
	activity, err := fit.Activity()
	if err != nil {
		return err
	}

	// Print the latitude and longitude of the first Record message
	for _, record := range activity.Records {
		fmt.Println(record.PositionLat)
		fmt.Println(record.PositionLong)
		break
	}

	// Print the sport of the first Session message
	for _, session := range activity.Sessions {
		fmt.Println(session.Sport)
		break
	}

	// Output:
	// 2012-04-09 21:22:26 +0000 UTC
	// Hrm1
	// Activity
	// 41.51393
	// -73.14859
	// Running
	return nil
}

func init() {
	rootCmd.AddCommand(fitCmd)
}

var fitCmd = &cobra.Command{
	Use:   "fit",
	Short: "Run fit",
	Long:  `Run fit`,
	// Aliases: []string{"g"},
	RunE: runfit,
}
