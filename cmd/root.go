package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/kvannotten/pcd"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "pcd",
	Short: "CLI podcatcher (podcast client)",
	Long: `pcd is a CLI application that allows you to track and download your podcasts.
Just add the necessary configuration under ~/.config/pcd and you can get started. 
Run pcd -h to get full help.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/pcd)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".pcd" (without extension).
		viper.AddConfigPath(filepath.Join(home, ".config"))
		viper.SetConfigName("pcd")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func findPodcast(idOrName string) (*pcd.Podcast, error) {
	id, err := strconv.Atoi(idOrName)
	if err != nil {
		switch err.(type) {
		case *strconv.NumError: // try as a name instead
			return findByName(idOrName), nil
		default:
			log.Print("Could not parse podcast search argument, please use the ID or the name")
			return nil, err
		}
	}
	return findByID(id), nil
}

func findByName(name string) *pcd.Podcast {
	return findByFunc(func(podcast *pcd.Podcast) bool {
		return podcast.Name == name
	})
}

func findByID(id int) *pcd.Podcast {
	return findByFunc(func(podcast *pcd.Podcast) bool {
		return podcast.ID == id
	})
}

func findByFunc(fn func(podcast *pcd.Podcast) bool) *pcd.Podcast {
	var podcasts []pcd.Podcast

	if err := viper.UnmarshalKey("podcasts", &podcasts); err != nil {
		log.Fatalf("Could not parse 'podcasts' entry in config: %v", err)
	}

	for _, podcast := range podcasts {
		if fn(&podcast) {
			return &podcast
		}
	}

	return nil
}
