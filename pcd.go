package main

import (
	"fmt"
	"github.com/kvannotten/pcd/configuration"
	"github.com/kvannotten/pcd/feedparser"
	"os"
	"strconv"

	"github.com/codegangsta/cli"
	"github.com/dustin/go-humanize"
)

var (
	conf *configuration.Config
)

func main() {
	conf = configuration.InitConfiguration()
	err := conf.Validate()
	if err != nil {
		panic(err)
	}

	app := cli.NewApp()
	app.Name = "pcd"
	app.Usage = "CLI podcast client"
	app.Version = "0.0.1"

	app.Commands = []cli.Command{
		{
			Name:    "sync",
			Aliases: []string{"s"},
			Usage:   "Sync podcasts defined in ~/.pcd",
			Action:  SyncPodcasts,
		},
		{
			Name:    "download",
			Aliases: []string{"d"},
			Usage:   "Download podcast episode: `pcd download <podcast_id/name> <episode_offset>`. The <episode_offset> is the chronological number of the episode where 1 is the latest.",
			Action:  DownloadPodcast,
		},
		{
			Name:    "play",
			Aliases: []string{"p"},
			Usage:   "Play specified podcast",
			Action:  PlayPodcast,
		},
		{
			Name:    "list",
			Aliases: []string{"l"},
			Usage:   "List the episodes of a specified podcast: `pcd list <podcast_id/name>`",
			Action:  ListPodcast,
		},
	}

	app.Run(os.Args)
}

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
	fmt.Printf("\t%-20s - %6s - %10s\n\t---------------------------------------------------\n", "Name", "Size", "Downloaded?")
	for _, item := range items {
		fmt.Printf("\t%-20s - %6s - %t\n", item.Title.Title, humanize.Bytes(item.Enclosure.Length), item.Downloaded)
	}
	fmt.Println()

}
