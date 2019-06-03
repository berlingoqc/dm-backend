package tasks

import (
	"github.com/berlingoqc/dm/tr"
	"github.com/berlingoqc/dm/file"
)

// CPTask ...
type CPTask struct {
	destination string
}

// Get ...
func (c *CPTask) Get() tr.ITask {
	return &CPTask{}
}

// GetID ...
func (c *CPTask) GetID() string {
	return "CPP"
}

// GetInfo ...
func (c *CPTask) GetInfo() tr.TaskInfo {
	return tr.TaskInfo{
		Name:        "copy",
		Description: "Copy file to another location , work on file and directory",
		Params: []tr.Params{
			tr.Params{
				Name:        "destination",
				Type:        "string",
				Optional:    false,
				Description: "Folder where to copy the file",
			},
		},
	}
}

// Execute ...
func (c *CPTask) Execute(filepath string, params map[string]interface{}, channel chan tr.TaskFeedBack) {
	var err error
	defer func() {
		if err != nil {
			channel <- tr.TaskFeedBack{
				Event: "ERROR",
				Message: err,
			}
		}
	}()
	c.destination = params["destination"].(string)
	println("CP : copying ",filepath, " TO ", c.destination)
	if err = file.Copy(filepath,c.destination); err != nil {
		return
	}
	channel <- tr.TaskFeedBack{
		Event: "DONE",
		Message: tr.TaskOver{
			Files: []string{""},
		},
	}

}
