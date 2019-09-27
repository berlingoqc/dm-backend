package triggers

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
			// Regarde quoi faire pour ce trigger
		}
	}
}

func init() {
	go fileHandlerMainLoop()
}
