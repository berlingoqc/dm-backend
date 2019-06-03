package main

import (
	"flag"

	"github.com/berlingoqc/dm/dm"
)

func main() {
	var configfile string

	flag.StringVar(&configfile, "config", "", "the configuration file")
	flag.Parse()

	ws, err := dm.Load(configfile)
	if err != nil {
		panic(err)
	}

	ws.Start()

}
