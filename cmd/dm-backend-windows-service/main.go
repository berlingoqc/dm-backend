package main

import (
	"syscall"
	"time"

	"github.com/berlingoqc/dm-backend/program/service/windows"
	"github.com/berlingoqc/dm-backend/webserver"
	"github.com/getlantern/systray"
)

func main() {
	var ws *webserver.WebServer

	windows.ServiceStartFunc = func() {
		go func() {
			println("SLEEEPING BEFOIRE STARTING")
			time.Sleep(10 * time.Second)
			//ws = dm.Run()
		}()
	}

	windows.ServiceShutdownFunc = func() {
		ws.ChannelStop <- syscall.SIGINT
		systray.Quit()
	}

	windows.ServiceCommand("dm-backend", "Dm backend service")
}
