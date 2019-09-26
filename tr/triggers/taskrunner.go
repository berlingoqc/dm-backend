package triggers

import "github.com/berlingoqc/dm-backend/tr/pipeline"

// TriggerEvent ...
type TriggerEvent struct {
	// The name of the event that trigger the pipeline
	Event string
	// File on witch to execute the pipeline
	File  string
}

// RegisterTrigger ...
type RegisterTrigger struct {

}

// TriggerEventChannel ...
var TriggerEventChannel = make(chan TriggerEvent)

func fileHandlerMainLoop() {
	for {
		select {
		case trigger := <-TriggerEventChannel:
			println("EVENT ", trigger.Event, " for file ", trigger.File)
			if _, err := pipeline.StartFromRegister(trigger.File); err != nil {
				println("Error starting register pipeline ", err.Error())
			}
		}
	}
}

func init() {
	go fileHandlerMainLoop()
}
