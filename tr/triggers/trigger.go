package triggers

// Settings ...
type Settings struct {
	RemoveAfterRun bool              `json:"remove_after_run"`
	PipelineID     string            `json:"pipeline_id"`
	Data           map[string]string `json:"data"`
}

// WatchInfo ...
type WatchInfo struct {
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
	GetWatchInfo() *map[int64]WatchInfo
	Init(ch chan PipelineTrigger)
}

// Triggers ...
var Triggers = make(map[string]ITrigger)

// InitTriggers ...
func InitTriggers() chan PipelineTrigger {
	ch := make(chan PipelineTrigger, 5)
	for _, t := range Triggers {
		t.Init(ch)
	}
	return ch
}
