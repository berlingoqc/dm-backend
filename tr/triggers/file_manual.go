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
func (m *ManualFileTrigger) Init(triggerChannel chan PipelineTrigger) {
	m.ch = triggerChannel
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
