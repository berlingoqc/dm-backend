package triggers

import (
	"encoding/json"
)

// AddEvent ...
type AddEvent struct {
	Trigger string    `json:"trigger"`
	Info    WatchInfo `json:"info"`
}

// RPC ...
type RPC struct{}

// AddEvent ...
func (t *RPC) AddEvent(data interface{}) int64 {
	dta, _ := json.Marshal(data)
	var event AddEvent
	json.Unmarshal(dta, &event)
	if t, ok := Triggers[event.Trigger]; ok {
		i, err := t.AddWatch(event.Info.Event, event.Info.Param, event.Info.Settings)
		if err != nil {
			panic(err)
		}
		return i
	}
	panic("NO TRIGGER")
}
