package tasks

import (
	"path/filepath"

	"github.com/berlingoqc/dm/file"
	"github.com/berlingoqc/dm/tr"
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
func (c *CPTask) Execute(fp string, params map[string]interface{}, channel chan tr.TaskFeedBack) {
	var err error
	defer tr.SendError(channel, err)
	c.destination = params["destination"].(string)
	println("CP : copying ", fp, " TO ", c.destination)

	// Doit ajouter le nom du fichier a la fin du path de destination pour que ca marche
	fileName := filepath.Base(fp)
	c.destination = filepath.Join(c.destination, fileName)

	if err = file.Copy(fp, c.destination); err != nil {
		return
	}
	tr.SendDone(channel, tr.TaskOver{
		Files: []string{c.destination},
	})

}
