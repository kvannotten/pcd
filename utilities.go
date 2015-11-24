package main

import (
	"fmt"
	"os"
	"os/exec"
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
	var wg sync.WaitGroup
	for _, podcast := range conf.Podcasts {
		wg.Add(1)
		fmt.Printf("Checking '%s'...\n", podcast.Name)
		go feedparser.Parse(podcast, &wg)
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
		if err != nil {
			fmt.Printf("Cannot find podcast %s", c.Args()[1])
			return
		}
	}
	feedparser.Download(*podcast, number)
}

func PlayPodcast(c *cli.Context) {
	podcastID := c.Args().First()
	podcast := findPodcast(podcastID)

	filename := feedparser.GetFileNameForPodcastAndEpisode(*podcast, 1)

	out, err := exec.Command(conf.Commands.Player, filepath.Join(podcast.Path, filename)).Output()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(out)
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
