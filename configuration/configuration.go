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
package configuration

import (
	"io/ioutil"
	"log"
	"os/user"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Podcasts Podcasts
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
