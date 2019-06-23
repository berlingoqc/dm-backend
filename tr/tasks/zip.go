package tasks

import (
	"github.com/berlingoqc/dm-backend/file"
	"github.com/berlingoqc/dm-backend/tr"
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
func (c *ZipTask) Get() tr.ITask {
	return &ZipTask{}
}

// GetID ...
func (c *ZipTask) GetID() string {
	return "ZIP"
}

// GetInfo ...
func (c *ZipTask) GetInfo() tr.TaskInfo {
	return tr.TaskInfo{
		Name:        "zip",
		Description: "Unzip the archive to a directory",
		Params: []tr.Params{
			tr.Params{
				Name:        "methode",
				Type:        "string",
				Optional:    false,
				Description: "zip or unzip",
			},
			tr.Params{
				Name:        "destination",
				Type:        "string",
				Optional:    false,
				Description: "Folder where to unzip the archive",
			},
		},
	}
}

// Execute ...
func (c *ZipTask) Execute(filepath string, params map[string]interface{}, channel chan tr.TaskFeedBack) {
	var err error
	defer tr.SendError(channel, err)
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

	tr.SendDone(channel, tr.TaskOver{
		Files: []string{zipParam.Destination},
	})
}
