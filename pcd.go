package main

import (
	"fmt"
	"github.com/kvannotten/pcd/configuration"
)

var (
	Conf *configuration.Config
)

func main() {
	Conf = configuration.InitConfiguration()
	fmt.Printf("%#v", Conf)
}
