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
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:     "list <podcast_id/podcast_name>",
	Aliases: []string{"ls"},
	Short:   "Lists all episodes of a podcast",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 1 {
			podcast, err := findPodcast(args[0])
			if err != nil {
				log.Fatal("Could not perform search")
			}
			if podcast == nil {
				log.Fatalf("Could not find podcast with search: %s", args[0])
			}

			if err := podcast.Load(); err != nil {
				log.Fatalf("Could not load podcast: %#v", err)
			}

			fmt.Print(podcast)
		} else if len(args) == 0 {
			fmt.Println("List of podcasts from your configuration:")
			for _, podcast := range findAll() {
				if err := podcast.Load(); err != nil {
					log.Fatalf("Could not load podcast: %#v", err)
				}
                                fmt.Printf("\t%d - %-40s (%d episodes)\n", podcast.ID, podcast.Name, len(podcast.Episodes))
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
