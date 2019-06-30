package tasks

import (
	"fmt"
	"path/filepath"

	"github.com/berlingoqc/dm-backend/file"
	"github.com/berlingoqc/dm-backend/tr/task"
)

// CPTask ...
type CPTask struct {
	destination string
}

// Get ...
func (c *CPTask) Get() task.ITask {
	return &CPTask{}
}

// GetID ...
func (c *CPTask) GetID() string {
	return "CPP"
}

// GetInfo ...
func (c *CPTask) GetInfo() task.TaskInfo {
	return task.TaskInfo{
		Name:        "copy",
		Description: "Copy file to another location , work on file and directory",
		Params: []task.Params{
			task.Params{
				Name:        "destination",
				Type:        "string",
				Optional:    false,
				Description: "Folder where to copy the file",
			},
		},
		NumberReturn: 1,
	}
}

// Execute ...
func (c *CPTask) Execute(fp string, params map[string]interface{}, channel chan task.TaskFeedBack) {
	var err error
	var ok bool
	defer task.SendError(channel, err)
	if c.destination, ok = params["destination"].(string); !ok {
		channel <- task.TaskFeedBack{Event: task.ErrorFeedBack, Message: "destination param not valid"}
	}
	output := fmt.Sprint("CP : copying  TO ", c.destination)
	channel <- task.TaskFeedBack{Event: task.OutFeedBack, Message: output}

	// Doit ajouter le nom du fichier a la fin du path de destination pour que ca marche
	fileName := filepath.Base(fp)
	c.destination = filepath.Join(c.destination, fileName)

	if err = file.Copy(fp, c.destination); err != nil {
		channel <- task.TaskFeedBack{Event: task.ErrorFeedBack, Message: err}
		return
	}
	task.SendDone(channel, task.TaskOver{
		Files: []string{c.destination},
	})
}
