package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/bzimmer/gravl/pkg/common"
	"github.com/mitchellh/go-homedir"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	// The name of our config file without the file extension
	defaultConfigFilename = ".gravl"

	// The environment variable prefix of all environment variables bound to our command line flags.
	// For example, --number is bound to GRAVL_NUMBER.
	envPrefix = "GRAVL"
)

var (
	debug      bool
	compact    bool
	monochrome bool
	verbosity  string
	config     string
	encoder    *json.Encoder
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gravl",
	Short: "Tools for planning adventures",
	Long:  `Planning outdoor adventures since 2020.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if err := initLogging(cmd); err != nil {
			return nil
		}
		if err := initConfig(cmd); err != nil {
			return nil
		}
		encoder = common.NewEncoder(cmd.OutOrStdout(), compact)
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	ctx := context.Background()
	if err := rootCmd.ExecuteContext(ctx); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(
		&config, "config", "", "config file (default is $HOME/.gravl.yaml)")
	rootCmd.PersistentFlags().StringVarP(&verbosity, "verbosity", "v",
		zerolog.InfoLevel.String(), "Log level (trace, debug, info, warn, error, fatal, panic")
	rootCmd.PersistentFlags().BoolVarP(&monochrome, "monochrome", "m",
		false, "Use monochrome logging, color enabled by default")
	rootCmd.PersistentFlags().BoolVarP(&compact, "compact", "c",
		false, "Use compact JSON output")
}

func initLogging(cmd *cobra.Command) error {
	level, err := zerolog.ParseLevel(verbosity)
	if err != nil {
		return err
	}
	debug = (level == zerolog.DebugLevel)
	zerolog.SetGlobalLevel(level)
	zerolog.DurationFieldUnit = time.Millisecond
	zerolog.DurationFieldInteger = true
	log.Logger = log.Output(
		zerolog.ConsoleWriter{
			Out:     cmd.OutOrStderr(),
			NoColor: monochrome,
		},
	)
	return nil
}

// Shamelessly stolen from:
//  https://carolynvanslyck.com/blog/2020/08/sting-of-the-viper/
//
// Order of precedence: default < config file < env variable < cli flag
func initConfig(cmd *cobra.Command) error {
	v := viper.New()

	if config != "" {
		// Use config file from the flag
		v.SetConfigFile(config)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			return err
		}
		// Set the config path to be the user's home directory
		v.AddConfigPath(home)
		// Set the base name of the config file, without the file extension
		v.SetConfigName(defaultConfigFilename)
	}

	// Attempt to read the config file, gracefully ignoring errors
	// caused by a config file not being found. Return an error
	// if we cannot parse the config file.
	if err := v.ReadInConfig(); err != nil {
		log.Debug().Str("path", v.ConfigFileUsed()).Msg("config")
		// It's okay if there isn't a config file
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Warn().Err(err).Str("path", v.ConfigFileUsed()).Msg("config")
			return err
		}
	}

	// When we bind flags to environment variables expect that the
	// environment variables are prefixed, e.g. a flag like --number
	// binds to an environment variable GRAVL_NUMBER. This helps
	// avoid conflicts.
	v.SetEnvPrefix(envPrefix)

	// Bind to environment variables
	// Works great for simple config names, but needs help for names
	// like --favorite-color which we fix in the bindFlags function
	v.AutomaticEnv()

	// Bind the current command's flags to viper
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		// Environment variables can't have dashes in them, so bind them to their
		// equivalent keys with underscores, e.g. --favorite-color to GRAVL_FAVORITE_COLOR
		if strings.Contains(f.Name, "-") {
			envVarSuffix := strings.ToUpper(strings.ReplaceAll(f.Name, "-", "_"))
			v.BindEnv(f.Name, fmt.Sprintf("%s_%s", envPrefix, envVarSuffix))
		}

		// Apply the viper config value to the flag when the flag is not set and viper has a value
		if !f.Changed && v.IsSet(f.Name) {
			val := v.Get(f.Name)
			cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
		}
	})
	return nil
}
