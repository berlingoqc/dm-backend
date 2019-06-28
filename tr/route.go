package tr

// RPC ...
type RPC struct{}

// TriggerRegister ...
func (t *RPC) TriggerRegister(event, file string) {
	TriggerEventChannel <- TriggerEvent{
		Event: event,
		File:  file,
	}
}
