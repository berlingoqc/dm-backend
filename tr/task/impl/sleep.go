package tasks

import (
	"strconv"
	"time"

	"github.com/berlingoqc/dm-backend/tr/task"
)

// SleepTask ...
type SleepTask struct {
}

// Get ...
func (s *SleepTask) Get() task.ITask {
	return &SleepTask{}
}

// GetID ...
func (s *SleepTask) GetID() string {
	return s.GetInfo().Name
}

// GetInfo ...
func (s *SleepTask) GetInfo() task.TaskInfo {
	return task.TaskInfo{
		Name:        "sleep",
		Provider:    "native",
		Description: "Sleep for a amount of time and return the same file",
		Params: []task.Params{
			task.Params{
				Name:        "duration",
				Type:        "number",
				Optional:    false,
				Description: "Number of millisecond of sleep",
			},
		},
		Return: []task.Return{
			task.Return{
				Description: "The same file",
			},
		},
	}
}

// Execute ...
func (s *SleepTask) Execute(file string, params map[string]interface{}, channel chan task.TaskFeedBack) {
	d, ok := params["duration"].(string)
	if !ok {
		channel <- task.TaskFeedBack{Event: task.ErrorFeedBack, Message: "invalid duration"}
		return
	}
	duration, _ := strconv.Atoi(d)
	time.Sleep(time.Duration(duration) * time.Millisecond)

	task.SendDone(channel, task.TaskOver{
		Files: []string{file},
	})
}
