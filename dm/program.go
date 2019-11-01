package dm

import (
	"flag"
	"path"

	"github.com/berlingoqc/dm-backend/webserver"
	"github.com/mitchellh/go-homedir"
)

var ws *webserver.WebServer

// Run ...
func Run() *webserver.WebServer {

	var configfile string

	flag.StringVar(&configfile, "config", "", "the configuration file")
	flag.Parse()

	if configfile == "" {
		dir, _ := homedir.Dir()
		configfile = path.Join(dir, ".dm", "config.json")
	}

	var err error
	ws, err = Load(configfile)
	if err != nil {
		panic(err)
	}
	ws.StartAsync()
	return ws
}
