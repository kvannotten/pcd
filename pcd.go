package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os/user"

	"gopkg.in/yaml.v2"
)

var (
	Conf Config
)

func main() {
	InitConfiguration()
	fmt.Printf("%#v", Conf)
}

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
	Username string
	Password string
}

type Podcasts []Podcast

func InitConfiguration() {
	usr, err := user.Current()
	if err != nil {
		log.Println("Could not find user.")
	}
	homeDir := usr.HomeDir

	source, err := ioutil.ReadFile(homeDir + "/.pcd")
	if err != nil {
		log.Println("Could not read configuration file.")
	}

	err = yaml.Unmarshal(source, &Conf)
	if err != nil {
		log.Println("Could not parse configuration file.")
	}
}
