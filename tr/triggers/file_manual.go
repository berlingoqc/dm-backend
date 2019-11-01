package triggers

// ManualFileTrigger ...
type ManualFileTrigger struct {
	ch chan PipelineTrigger
}

// GetWatchInfo ...
func (m *ManualFileTrigger) GetWatchInfo() *map[int64]WatchInfo {
	return nil
}

// Init ...
func (m *ManualFileTrigger) Init(triggerChannel chan PipelineTrigger, _ chan interface{}) {
	m.ch = triggerChannel
}

// DeleteWatch ...
func (m *ManualFileTrigger) DeleteWatch(id int64) error {
	return nil
}

// AddWatch ...
func (m *ManualFileTrigger) AddWatch(event string, param interface{}, settings *Settings) (int64, error) {
	m.ch <- PipelineTrigger{
		File:       param.(string),
		PipelineID: settings.PipelineID,
		Data:       settings.Data,
	}
	return 0, nil
}
