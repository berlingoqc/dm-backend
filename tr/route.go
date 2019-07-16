package tr

// RPC ...
type RPC struct{}

// TriggerRegister to manually trigger a register task
func (t *RPC) TriggerRegister(event, file string) {
	TriggerEventChannel <- TriggerEvent{
		Event: event,
		File:  file,
	}
}
