/*
	pcd - Simple, lightweight podcatcher in golang
	Copyright (C) 2016  Kristof Vannotten

	This program is free software: you can redistribute it and/or modify
	it under the terms of the GNU General Public License as published by
	the Free Software Foundation, either version 3 of the License, or
	(at your option) any later version.

	This program is distributed in the hope that it will be useful,
	but WITHOUT ANY WARRANTY; without even the implied warranty of
	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
	GNU General Public License for more details.

	You should have received a copy of the GNU General Public License
	along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/codegangsta/cli"
	"github.com/dustin/go-humanize"
	"github.com/kvannotten/pcd/configuration"
	"github.com/kvannotten/pcd/feedparser"
	"github.com/olekukonko/tablewriter"
)

func SyncPodcasts(c *cli.Context) {
	var throttle = make(chan int, 4)
	var wg sync.WaitGroup

	for _, podcast := range conf.Podcasts {
		throttle <- 1
		wg.Add(1)
		fmt.Printf("Checking '%s' [id: %d]...\n", podcast.Name, podcast.ID)
		go feedparser.Parse(podcast, &wg, throttle)
	}
	wg.Wait()

	fmt.Printf("Done\n")
}

func DownloadPodcast(c *cli.Context) {
	if len(c.Args()) < 1 {
		fmt.Println("Please specify a podcast to download")
		return
	}

	podcast := findPodcast(c.Args().First())
	number := 1
	if len(c.Args()) > 1 {
		var err error
		number, err = strconv.Atoi(c.Args()[1])
		if number > len(feedparser.ListEpisodes(*podcast)) {
			fmt.Printf("There's not that many episodes")
			return
		}
		if err != nil {
			fmt.Printf("Cannot find podcast %s", c.Args()[1])
			return
		}
	}
	feedparser.Download(*podcast, number)
}

func ListPodcast(c *cli.Context) {
	podcastID := c.Args().First()
	if len(podcastID) > 0 {
		// list episodes of podcast
		printPodcastInfo(*findPodcast(podcastID))
	} else {
		// list all podcasts
		for _, podcast := range conf.Podcasts {
			printPodcastInfo(podcast)
		}
	}
}

func GetPodcastPath(c *cli.Context) {
	if len(c.Args()) < 1 {
		fmt.Println("Please specify a podcast id")
		return
	}

	podcast := findPodcast(c.Args().First())
	number := 1
	if len(c.Args()) > 1 {
		var err error
		number, err = strconv.Atoi(c.Args()[1])
		if err != nil {
			fmt.Printf("Cannot find podcast episode %s", c.Args()[1])
			return
		}
	}
	path := filepath.Join(podcast.Path, feedparser.GetFileNameForPodcastAndEpisode(*podcast, number))
	fmt.Println(path)
}

func findPodcast(searchTerm interface{}) *configuration.Podcast {
	id, _ := strconv.Atoi(searchTerm.(string))
	for _, podcast := range conf.Podcasts {
		if podcast.Name == searchTerm || podcast.ID == id {
			return &podcast
		}
	}
	return &configuration.Podcast{}
}

func printPodcastInfo(podcast configuration.Podcast) {
	fmt.Printf("ID: %d\n", podcast.ID)
	fmt.Printf("Name: %s\n", podcast.Name)
	fmt.Printf("Feed: %s\n", podcast.Feed)
	fmt.Printf("Path: %s\n", podcast.Path)
	fmt.Printf("Episodes: \n")

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Name", "Size", "Published Date", "Downloaded"})
	table.SetBorder(false)

	items := feedparser.ListEpisodes(podcast)
	for id, item := range items {
		var length uint64 = uint64(item.Enclosure.Length)
		var downloaded string
		if item.Downloaded {
			downloaded = "true"
		} else {
			downloaded = "false"
		}
		table.Append([]string{fmt.Sprintf("%d", id+1), item.Title.Title, humanize.Bytes(length), item.Date.Date, downloaded})
	}
	table.Render()
}
