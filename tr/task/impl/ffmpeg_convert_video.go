package tasks

import "github.com/berlingoqc/dm-backend/tr/task"

// FFMPEGConvertVideo ...
type FFMPEGConvertVideo struct {
}

// Get ...
func (f *FFMPEGConvertVideo) Get() task.ITask {
	return &FFMPEGConvertVideo{}
}

// GetID ...
func (f *FFMPEGConvertVideo) GetID() string {
	return "ffmpeg_convert_video"
}

// GetInfo ...
func (f *FFMPEGConvertVideo) GetInfo() task.TaskInfo {
	return task.TaskInfo{}
}

// Execute ...
func (f *FFMPEGConvertVideo) Execute(file string, params map[string]interface{}, channel chan task.TaskFeedBack) {
	panic("not implemented")
}
