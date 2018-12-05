// Copyright © 2018 Kristof Vannotten <kristof@vannotten.be>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package cmd

import (
	"log"

	"github.com/kvannotten/pcd"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// syncCmd represents the sync command
var syncCmd = &cobra.Command{
	Use:     "sync",
	Aliases: []string{"s"},
	Short:   "Syncs your podcasts",
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
