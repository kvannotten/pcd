package main

import (
	"github.com/kvannotten/pcd/configuration"
	"github.com/kvannotten/pcd/feedparser"
	"os"

	"github.com/codegangsta/cli"
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
			Usage:   "Download the podcast with the specified id/name",
			Action:  DownloadPodcast,
		},
		{
			Name:    "play",
			Aliases: []string{"p"},
			Usage:   "Play specified podcast",
			Action:  PlayPodcast,
		},
	}

	app.Run(os.Args)
}

func SyncPodcasts(c *cli.Context) {
	for i := 0; i < len(conf.Podcasts); i++ {
		feedparser.Parse(conf.Podcasts[i])
	}
}

func DownloadPodcast(c *cli.Context) {
	feedparser.Download(conf.Podcasts[0])
}

func PlayPodcast(c *cli.Context) {

}
