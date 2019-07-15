package main

import (
	"syscall"

	"github.com/berlingoqc/dm-backend/dm"
	"github.com/berlingoqc/dm-backend/webserver"
	"github.com/berlingoqc/dm-backend/program/service/windows"
	"github.com/getlantern/systray"
)

func main() {
	var ws *webserver.WebServer

	windows.ServiceStartFunc = func() {
		ws = dm.Run()
	}

	windows.ServiceShutdownFunc = func() {
		ws.ChannelStop <- syscall.SIGINT
		systray.Quit()
	}

	windows.ServiceCommand("dm-backend")
}
