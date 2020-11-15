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
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/cheggaaa/pb"
	"github.com/kvannotten/pcd"
	"github.com/spf13/cobra"
)

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:     "download <podcast> <episode_id>",
	Aliases: []string{"d"},
	Short:   "Downloads an episode of a podcast.",
	Long: `
This command will download one or multiple episode(s) of a podcast that you 
define. 

The episode number can be obtained by running 'pcd ls <podcast>' 

For example:

To download one episode

pcd ls gnu_open_world
pcd download gnu_open_world 1

To download episode ranges:

pcd download gnu_open_world '20-30,!25'

This will download episode 20 to 30 and skip the 25.

Available formats:

Episode numbers: '1,5,105'
Ranges: '2-15'
Skipping: '!102,!121'

Combining those as follow:

pcd download gnu_open_world '1-30,40-47,!15,!17,!20,102'

Make sure to use the single-quote on bash otherwise the !105 will expand your 
bash history.`,
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

	if len(args) < 2 {
		// download latest
		downloadEpisode(podcast, len(podcast.Episodes))
	}

	episodes, err := parseRangeArg(args[1])
	if err != nil {
		log.Fatalf("Could not parse episode number %s: %#v", args[1], err)
	}

	for _, n := range episodes {
		downloadEpisode(podcast, n)
	}
}

func downloadEpisode(podcast *pcd.Podcast, episodeN int) {
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

// parseRangeArg parses episodes number with the following format
// 1,2,3-5,!4 returns [1, 2, 3, 5]
func parseRangeArg(arg string) ([]int, error) {
	if len(arg) == 0 {
		return nil, nil
	}

	// we try to convert the arg as a single episode
	n, err := strconv.Atoi(arg)
	if err == nil {
		return []int{n}, nil
	}

	var results []int

	// extract negative numbers !X
	negatives := regexp.MustCompile(`!\d+`)
	notWanted := negatives.FindAllString(arg, -1)

	arg = negatives.ReplaceAllString(arg, "")

	// extract ranges X-Y
	rangesPattern := regexp.MustCompile(`\d+-\d+`)
	ranges := rangesPattern.FindAllString(arg, -1)

	arg = rangesPattern.ReplaceAllString(arg, "")

	// extract the remaining single digit X
	digitsPattern := regexp.MustCompile(`\d+`)
	digits := digitsPattern.FindAllString(arg, -1)

	for _, r := range ranges {
		parts := strings.Split(r, "-")
		if len(parts) != 2 {
			return nil, fmt.Errorf("range %s must have the format start-end", r)
		}

		start, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, err
		}

		end, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, err
		}

		for i := start; i <= end; i++ {
			// make sure it's wanted
			wanted := true
			for _, nw := range notWanted {
				if fmt.Sprintf("!%d", i) == nw {
					wanted = false
					break
				}
			}

			if !wanted {
				continue
			}

			results = append(results, i)
		}
	}

	// let's add the remaining digits
	for _, d := range digits {
		i, err := strconv.Atoi(d)
		if err != nil {
			return nil, err
		}

		results = append(results, i)
	}

	// we sort the result
	sort.Ints(results)

	return results, nil
}
