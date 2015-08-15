package main

import (
	"fmt"
	"strconv"

	"github.com/codegangsta/cli"
	"github.com/dustin/go-humanize"
	"github.com/kvannotten/pcd/configuration"
	"github.com/kvannotten/pcd/feedparser"
)

func SyncPodcasts(c *cli.Context) {
	for _, podcast := range conf.Podcasts {
		fmt.Printf("Checking '%s'...", podcast.Name)
		feedparser.Parse(podcast)
		fmt.Printf(" Done\n")
	}
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
		if err != nil {
			fmt.Printf("Cannot find podcast %s", c.Args()[1])
			return
		}
	}
	feedparser.Download(*podcast, number)
}

func PlayPodcast(c *cli.Context) {

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
	fmt.Printf("Name: %s\n", podcast.Name)
	fmt.Printf("Feed: %s\n", podcast.Feed)
	fmt.Printf("Path: %s\n", podcast.Path)
	items := feedparser.ListEpisodes(podcast)
	fmt.Printf("Episodes: \n")
	fmt.Printf("\t%4s - %-25s - %6s - %10s\n\t---------------------------------------------------\n", "ID", "Name", "Size", "Downloaded?")
	for id, item := range items {
		var title string
		if len(item.Title.Title) > 25 {
			title = item.Title.Title[0:25]
		} else {
			title = item.Title.Title
		}
		fmt.Printf("\t%4d - %-25s - %6s - %t\n", id+1, title, humanize.Bytes(item.Enclosure.Length), item.Downloaded)
	}
	fmt.Println()
}
