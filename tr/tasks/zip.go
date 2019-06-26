package tasks

import (
	"github.com/berlingoqc/dm-backend/file"
	"github.com/berlingoqc/dm-backend/tr/task"
	"github.com/mitchellh/mapstructure"
)

// ZipParams ...
type ZipParams struct {
	Methode     string `json:"methode"`
	Destination string `json:"destination"`
}

// ZipTask ...
type ZipTask struct {
}

// Get ...
func (c *ZipTask) Get() task.ITask {
	return &ZipTask{}
}

// GetID ...
func (c *ZipTask) GetID() string {
	return "ZIP"
}

// GetInfo ...
func (c *ZipTask) GetInfo() task.TaskInfo {
	return task.TaskInfo{
		Name:        "zip",
		Description: "Unzip the archive to a directory",
		Params: []task.Params{
			task.Params{
				Name:        "methode",
				Type:        "string",
				Optional:    false,
				Description: "zip or unzip",
			},
			task.Params{
				Name:        "destination",
				Type:        "string",
				Optional:    false,
				Description: "Folder where to unzip the archive",
			},
		},
	}
}

// Execute ...
func (c *ZipTask) Execute(filepath string, params map[string]interface{}, channel chan task.TaskFeedBack) {
	var err error
	defer task.SendError(channel, err)
	zipParam := ZipParams{}
	if err = mapstructure.Decode(params, &zipParam); err != nil {
		return
	}
	switch zipParam.Methode {
	case "zip":
		println("NOT SUPPORTED ZIP TEY")
	case "unzip":
		_, err = file.Unzip(filepath, zipParam.Destination)
		if err != nil {
			return
		}
	}

	task.SendDone(channel, task.TaskOver{
		Files: []string{zipParam.Destination},
	})
}
