package triggers

import (
	"math/rand"
	"time"
)

// Settings ...
type Settings struct {
	RemoveAfterRun bool              `json:"remove_after_run"`
	PipelineID     string            `json:"pipeline_id"`
	Data           map[string]string `json:"data"`
}

// WatchInfo ...
type WatchInfo struct {
	Trigger  string      `json:"trigger"`
	Event    string      `json:"event"`
	Param    interface{} `json:"param"`
	Settings *Settings   `json:"settings"`
}

// PipelineTrigger ...
type PipelineTrigger struct {
	File       string            `json:"file"`
	PipelineID string            `json:"pipeline_id"`
	Data       map[string]string `json:"data"`
}

// ITrigger ...
type ITrigger interface {
	AddWatch(event string, param interface{}, settings *Settings) (int64, error)
	DeleteWatch(id int64) error
	GetWatchInfo() *map[int64]WatchInfo
	Init(ch chan PipelineTrigger, signal chan interface{})
}

// Triggers ...
var Triggers = make(map[string]ITrigger)

// InitTriggers ...
func InitTriggers() (chan PipelineTrigger, chan interface{}) {
	ch := make(chan PipelineTrigger, 5)
	chSignal := make(chan interface{}, 5)
	for _, t := range Triggers {
		t.Init(ch, chSignal)
	}
	return ch, chSignal
}

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

func getID() int64 {
	return r.Int63()
}
