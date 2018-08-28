package cmd

import (
	"log"

	"github.com/kvannotten/pcd"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// syncCmd represents the sync command
var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Syncs your podcasts",
	Run: func(cmd *cobra.Command, args []string) {
		var podcasts []pcd.Podcast

		if err := viper.UnmarshalKey("podcasts", &podcasts); err != nil {
			log.Fatalf("Could not parse 'podcasts' entry in config: %v", err)
		}

		for _, podcast := range podcasts {
			log.Printf("[%s] Syncing...", podcast.Name)
			if err := podcast.Sync(); err != nil {
				log.Printf("[%s] Could not sync podcast: %v", podcast.Name, err)
				continue
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
}
