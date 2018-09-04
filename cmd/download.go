package cmd

import (
	"log"
	"strconv"

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
	Args: cobra.MinimumNArgs(2),
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

	episodeN, err := strconv.Atoi(args[1])
	if err != nil {
		log.Fatalf("Could not parse episode number %s: %#v", args[1], err)
	}

	if episodeN > len(podcast.Episodes) {
		log.Fatalf("There's only %d episodes in this podcast.", len(podcast.Episodes))
	}

	if episodeN < 1 {
		log.Fatalf("A number from 1 to %d is required.", len(podcast.Episodes))
	}

	bar := pb.New(podcast.Episodes[episodeN-1].Length).SetUnits(pb.U_BYTES)
	bar.ShowTimeLeft = true
	bar.ShowSpeed = true
	bar.Start()
	if err := podcast.Episodes[episodeN-1].Download(podcast.Path, bar); err != nil {
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
