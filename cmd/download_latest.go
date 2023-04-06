package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var downloadLatestCmd = &cobra.Command{
	Use:     "download-latest",
	Aliases: []string{"dl"},
	Short:   "Downloads the latest episode from each podcast feed",
	Long: `
This command will download the latest episode from each podcast

For example

pcd da
`,
	Run: downloadLatest,
}

func init() {
	rootCmd.AddCommand(downloadLatestCmd)
}

func downloadLatest(cmd *cobra.Command, args []string) {
	podcasts := findAll()

	for _, singlePodcast := range podcasts {

		if err := singlePodcast.Load(); err != nil {
			log.Println("Could not load podcast: ", err)
			continue
		}

		if err := downloadEpisode(&singlePodcast, len(singlePodcast.Episodes)); err != nil {
			log.Println(err)
		}
	}
}
