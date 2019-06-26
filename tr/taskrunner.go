package tr

import "github.com/berlingoqc/dm-backend/tr/pipeline"

// TriggerEvent ...
type TriggerEvent struct {
	Event string
	File  string
}

// TriggerEventChannel ...
var TriggerEventChannel = make(chan TriggerEvent)

func fileHandlerMainLoop() {
	for {
		select {
		case trigger := <-TriggerEventChannel:
			println("EVENT ", trigger.Event, " for file ", trigger.File)
			go pipeline.Start(trigger.File)
		}
	}

}

func init() {
	go fileHandlerMainLoop()
}
