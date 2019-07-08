package dm

import (
	"flag"
	"path"
	"syscall"

	"github.com/berlingoqc/dm-backend/webserver"
	"github.com/getlantern/systray"
	"github.com/getlantern/systray/example/icon"
	"github.com/mitchellh/go-homedir"
)

var ws *webserver.WebServer

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
	go initSystray()
	return ws
}

func initSystray() {
	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetIcon(icon.Data)
	systray.SetTitle("DM Backend")
	systray.SetTooltip("DM Backend listening on ")
	mQuit := systray.AddMenuItem("Quit", "Quit the whole app")
	go func() {
		<-mQuit.ClickedCh
		ws.ChannelStop <- syscall.SIGINT
		systray.Quit()
	}()
	// Sets the icon of a menu item. Only available on Mac.
	mQuit.SetIcon(icon.Data)
}

func onExit() {

}
