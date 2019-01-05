package cmd

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	version   = "devel"
	commit    = "unknown"
	buildDate = "unknown"
	cfgFile   string
)

var rootCmd = &cobra.Command{
	Use:          "bookmarks",
	Short:        "Personal zero-touch bookmarking app in the cloud, with full text search support.",
	SilenceUsage: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if viper.GetBool("debug") {
			zerolog.SetGlobalLevel(zerolog.DebugLevel)
		} else {
			zerolog.SetGlobalLevel(zerolog.InfoLevel)
		}

		return nil
	},
}

// Execute executes the rootCmd logic and is the main entry point
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is .bookmarks.yaml in $PWD, $HOME, /etc)")

	rootCmd.PersistentFlags().BoolP("debug", "d", false, "Enable debug mode")
	rootCmd.PersistentFlags().StringP("storage", "s", "", "The location where to store state")

	viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))
	viper.BindPFlag("storage", rootCmd.PersistentFlags().Lookup("storage"))
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName(".bookmarks")
		viper.AddConfigPath(".")
		viper.AddConfigPath("$HOME")
		viper.AddConfigPath("/etc/")
	}

	viper.SetEnvPrefix("bookmarks")
	viper.AutomaticEnv()

	viper.ReadInConfig()
}
