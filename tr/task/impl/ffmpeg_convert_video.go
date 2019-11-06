package tasks

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/berlingoqc/dm-backend/tr/task"
)

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
	return task.TaskInfo{
		Name:        "ffmpeg_convert_video",
		Provider:    "native",
		Description: "Convert all video of unwanted format to a better format to be send over network",
		Params: []task.Params{
			task.Params{
				Name:         "ouput_format",
				Type:         "string",
				DefaultValue: "mp4",
				Optional:     false,
			},
			task.Params{
				Name:         "input_format",
				Type:         "[]string",
				DefaultValue: `["avi","wav"]`,
				Optional:     false,
			},
		},
	}
}

// Execute ...
func (f *FFMPEGConvertVideo) Execute(file string, params map[string]interface{}, channel chan task.TaskFeedBack) {
	ouputFormat := params["output_format"].(string)
	inputFormatStr := params["input_format"].(string)
	var inputFormat []string
	if err := json.Unmarshal([]byte(inputFormatStr), &inputFormat); err != nil {
		channel <- task.TaskFeedBack{Event: task.ErrorFeedBack, Message: err.Error()}
		return
	}

	files, err := getAllFileOfExtensionRecursive(file, inputFormat)
	if err != nil {
		channel <- task.TaskFeedBack{Event: task.ErrorFeedBack, Message: err.Error()}
		return
	}

	for _, file := range files {
		cmd := exec.Command("ffmpeg", file)
		if err := cmd.Run(); err != nil {
			channel <- task.TaskFeedBack{Event: task.ErrorFeedBack, Message: err.Error()}
			return
		}
	}

}

func getAllFileOfExtensionRecursive(folder string, extensions []string) ([]string, error) {
	var list []string
	return list, filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if isItemInList(filepath.Ext(path), extensions) {
			list = append(list, path)
		}
		return nil
	})
}

func isItemInList(item string, list []string) bool {
	for _, i := range list {
		if i == item {
			return true
		}
	}
	return false
}
