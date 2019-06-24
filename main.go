package main

import (
	"flag"

	"github.com/berlingoqc/dm-backend/dm"
	"github.com/getlantern/systray"
)

func main() {

	//systray.Run(onReady, onExit)

	var configfile string

	flag.StringVar(&configfile, "config", "", "the configuration file")
	flag.Parse()

	ws, err := dm.Load(configfile)
	if err != nil {
		panic(err)
	}

	ws.Start()

}

func onReady() {
	systray.SetTitle("Awesome App")
	systray.SetTooltip("Pretty awesome超级棒")
	systray.AddMenuItem("Quit", "Quit the whole app")
}

func onExit() {

}
