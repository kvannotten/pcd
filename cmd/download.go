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
	"strconv"
        "net/http"

	"github.com/cheggaaa/pb"
	"github.com/spf13/cobra"
)

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:     "download <podcast> <episode_id>",
	Aliases: []string{"d"},
	Short:   "Downloads an episode of a podcast.",
	Long: `This command will download an episode of a podcast that you define. The episode number can
be obtained by running 'pcd ls <podcast>' For example:

pcd ls gnu_open_world
pcd download gnu_open_world 1`,
	Args: cobra.MinimumNArgs(1),
	Run:  download,
}

func download(cmd *cobra.Command, args []string) {
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

	var episodeN int
	if len(args) > 1 {
		episodeN, err = strconv.Atoi(args[1])
	} else {
		episodeN = len(podcast.Episodes) // download latest
	}
	if err != nil {
		log.Fatalf("Could not parse episode number %s: %#v", args[1], err)
	}

	if episodeN > len(podcast.Episodes) {
		log.Fatalf("There's only %d episodes in this podcast.", len(podcast.Episodes))
	}

	if episodeN < 1 {
		log.Fatalf("A number from 1 to %d is required.", len(podcast.Episodes))
	}

	episodeToDownload := podcast.Episodes[episodeN-1]
	log.Printf("Started downloading: '%s' episode %d of %s", episodeToDownload.Title, episodeN, podcast.Name)

        // RSS Feeds cannot be trusted to accurately or consistently report the length
        // of the episode file. Instead, make a request for the header and use the
        // Content-Length property to get an accurate size.
         resp, err := http.Head(episodeToDownload.URL)
         if err != nil {
             log.Fatalf("Request failed: %s\nError: %#v", episodeToDownload.Title, err)
         }

         if resp.StatusCode != http.StatusOK {
             log.Fatalf("Request failed: %s\nError: %#v", episodeToDownload.Title, err)
         }
         size, _ := strconv.Atoi(resp.Header.Get("Content-Length"))

	bar := pb.New(size).SetUnits(pb.U_BYTES)
	bar.ShowTimeLeft = true
	bar.ShowSpeed = true
	bar.Start()

	if err := episodeToDownload.Download(podcast.Path, bar); err != nil {
		log.Fatalf("Could not download episode: %#v", err)
	}

	bar.Finish()
}

func init() {
	rootCmd.AddCommand(downloadCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// downloadCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// downloadCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
