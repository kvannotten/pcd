// Copyright © 2018 Kristof Vannotten <kristof@vannotten.be>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
        "strings"

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
			return findByNameFragment(idOrName), nil
		default:
			log.Print("Could not parse podcast search argument, please use the ID or the name")
			return nil, err
		}
	}
	return findByID(id), nil
}

func findAll() []pcd.Podcast {
	var podcasts []pcd.Podcast

	if err := viper.UnmarshalKey("podcasts", &podcasts); err != nil {
		log.Fatalf("Could not parse 'podcasts' entry in config: %v", err)
	}

	return podcasts
}

func findByNameFragment(name string) *pcd.Podcast {
	return findByFunc(func(podcast *pcd.Podcast) bool {
		return strings.Contains(
                    strings.ToLower(podcast.Name), strings.ToLower(name))
	})
}

func findByID(id int) *pcd.Podcast {
	return findByFunc(func(podcast *pcd.Podcast) bool {
		return podcast.ID == id
	})
}

func findByFunc(fn func(podcast *pcd.Podcast) bool) *pcd.Podcast {
	podcasts := findAll()
        var matchedPodcasts = make([]*pcd.Podcast, 0)

	for _, podcast := range podcasts {
		if fn(&podcast) {
                    matchedPodcast := podcast
                    matchedPodcasts = append(matchedPodcasts, &matchedPodcast)
		}
	}
        log.Print(matchedPodcasts)
        if len(matchedPodcasts) == 1 {
		return matchedPodcasts[0]
        } else {
                log.Fatalf("Provided search term matched too many podcasts: %v", matchedPodcasts)
        }

	return nil
}
