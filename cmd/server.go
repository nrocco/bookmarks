package cmd

import (
	"context"
	"os"

	"github.com/nrocco/bookmarks/api"
	"github.com/nrocco/bookmarks/scheduler"
	"github.com/nrocco/bookmarks/storage"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run the Bend web application and rest api",
	RunE: func(cmd *cobra.Command, args []string) error {
		logger := log.With().Int("pid", os.Getpid()).Logger()

		logger.Info().
			Bool("debug", viper.GetBool("debug")).
			Int("interval", viper.GetInt("interval")).
			Str("listen", viper.GetString("listen")).
			Str("storage", viper.GetString("storage")).
			Msg("Starting bookmarks")

		// Setup the database
		store, err := storage.New(context.Background(), viper.GetString("storage"))
		if err != nil {
			logger.Fatal().Err(err).Msg("Could not open the database")
		}
		logger.Info().Str("storage", viper.GetString("storage")).Msg("Store ready")

		// Setup the http server
		api := api.New(logger, store, viper.GetString("username"), viper.GetString("password"))
		logger.Info().Str("address", "http://"+viper.GetString("listen")).Msg("API ready")

		if viper.GetInt("interval") != 0 {
			scheduler.New(store, viper.GetInt("interval"))
		} else {
			logger.Info().Msg("Scheduler is disabled")
		}

		// Run the http server
		if err := api.ListenAndServe(viper.GetString("listen")); err != nil {
			logger.Warn().Err(err).Msg("Stopped the api server")
		}
		logger.Info().Msg("Stopping bookmarks")

		return nil
	},
}

func init() {
	serverCmd.PersistentFlags().StringP("listen", "l", "0.0.0.0:3000", "Address to listen for HTTP requests on")
	serverCmd.PersistentFlags().IntP("interval", "i", 15, "Fetch new feeds with this interval in minutes (0 to disable)")
	serverCmd.PersistentFlags().StringP("username", "u", "", "Username for authentication")
	serverCmd.PersistentFlags().StringP("password", "p", "", "Password for authentication")

	viper.BindPFlag("listen", serverCmd.PersistentFlags().Lookup("listen"))
	viper.BindPFlag("interval", serverCmd.PersistentFlags().Lookup("interval"))
	viper.BindPFlag("username", serverCmd.PersistentFlags().Lookup("username"))
	viper.BindPFlag("password", serverCmd.PersistentFlags().Lookup("password"))

	rootCmd.AddCommand(serverCmd)
}
