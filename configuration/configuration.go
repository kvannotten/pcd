package configuration

import (
	"io/ioutil"
	"log"
	"os/user"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Commands Commands
	Podcasts Podcasts
}

type Commands struct {
	Player   string
	Download string
}

type Podcast struct {
	ID       int
	Name     string
	Feed     string
	Path     string
	Username string
	Password string
}

type Podcasts []Podcast

func InitConfiguration() *Config {
	var config Config

	usr, err := user.Current()
	if err != nil {
		log.Println("Could not find user.")
	}
	homeDir := usr.HomeDir

	source, err := ioutil.ReadFile(homeDir + "/.pcd")
	if err != nil {
		log.Println("Could not read configuration file.")
	}

	err = yaml.Unmarshal(source, &config)
	if err != nil {
		log.Println("Could not parse configuration file.")
	}

	return &config
}

func (c *Config) Validate() error {
	// TODO: validate configuration file
	return nil
}
