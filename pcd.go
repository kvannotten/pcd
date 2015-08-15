package main

import (
	"os"

	"github.com/codegangsta/cli"
	"github.com/kvannotten/pcd/configuration"
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
			Usage:   "Download podcast episode: `pcd download <podcast_id/name> <episode_ID>`. The <episode_ID> is the chronological number of the episode where 1 is the latest.",
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
			Usage:   "List the episodes (and their episode_ID) of a specified podcast: `pcd list <podcast_id/name>`",
			Action:  ListPodcast,
		},
	}

	app.Run(os.Args)
}
