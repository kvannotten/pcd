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
			Name:    "list",
			Aliases: []string{"l"},
			Usage:   "List the episodes (and their episode_ID) of a specified podcast: `pcd list <podcast_id/name>`",
			Action:  ListPodcast,
		},
		{
			Name:    "get",
			Aliases: []string{"g"},
			Usage:   "Echo the path of the episode: `pcd get <podcast_id/name> [<episode_id>]`",
			Action:  GetPodcastPath,
		},
	}

	app.Run(os.Args)
}
