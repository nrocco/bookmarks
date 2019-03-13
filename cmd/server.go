package cmd

import (
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
		log.Info().
			Bool("debug", viper.GetBool("debug")).
			Int("interval", viper.GetInt("interval")).
			Str("listen", viper.GetString("listen")).
			Str("storage", viper.GetString("storage")).
			Msg("Starting bookmarks")

		// Setup the database
		store, err := storage.New(viper.GetString("storage"))
		if err != nil {
			log.Fatal().Err(err).Msg("Could not open the database")
		}

		// Setup the http server
		api := api.New(store, !viper.GetBool("noauth"))

		if viper.GetInt("interval") != 0 {
			// Setup the periodic scheduler
			scheduler.New(store, viper.GetInt("interval"))
		} else {
			log.Info().Msg("Scheduler is disabled")
		}

		// Run the http server
		if err := api.ListenAndServe(viper.GetString("listen")); err != nil {
			log.Warn().Err(err).Msg("Stopped the api server")
		}

		log.Info().Msg("Stopping bookmarks")

		return nil
	},
}

func init() {
	serverCmd.PersistentFlags().StringP("listen", "l", "0.0.0.0:3000", "Address to listen for HTTP requests on")
	serverCmd.PersistentFlags().IntP("interval", "i", 15, "Fetch new feeds with this interval in minutes (0 to disable)")
	serverCmd.PersistentFlags().BoolP("noauth", "n", false, "Disable authentication")

	viper.BindPFlag("listen", serverCmd.PersistentFlags().Lookup("listen"))
	viper.BindPFlag("interval", serverCmd.PersistentFlags().Lookup("interval"))
	viper.BindPFlag("noauth", serverCmd.PersistentFlags().Lookup("noauth"))

	rootCmd.AddCommand(serverCmd)
}
