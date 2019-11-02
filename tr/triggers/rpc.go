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
	i, err := GetTrigger(event.Trigger).AddWatch(event.Info.Event, event.Info.Param, event.Info.Settings)
	if err != nil {
		panic(err)
	}
	return i
}

// RemoveEvent ...
func (t *RPC) RemoveEvent(trigger string, id float64) int64 {
	i := int64(id)
	tr := GetTrigger(trigger)
	if err := tr.DeleteWatch(i); err != nil {
		panic(err)
	}
	return i
}

// GetAllEvents ...
func (t *RPC) GetAllEvents() map[int64]WatchInfo {
	m := make(map[int64]WatchInfo)
	for _, t := range Triggers {
		p := t.GetWatchInfo()
		for i, w := range *p {
			m[i] = w
		}
	}
	return m
}

// GetTrigger ...
func GetTrigger(name string) ITrigger {
	if t, ok := Triggers[name]; ok {
		return t
	}
	panic("NO TRIGGER NAMED " + name)
}
