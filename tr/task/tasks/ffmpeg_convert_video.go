package tasks

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

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
			task.Params{
				Name:         "delete",
				Type:         "bool",
				DefaultValue: "false",
				Optional:     true,
			},
			task.Params{
				Name:         "gpu",
				Type:         "bool",
				DefaultValue: "false",
				Optional:     true,
				Description:  "If true convert the video using nvea, requires nvidia gpu",
			},
		},
		Return: []task.Return{
			task.Return{
				Type:        "...string",
				Description: "All of the video file that were converted with there new extension",
			},
		},
	}
}

// Execute ...
func (f *FFMPEGConvertVideo) Execute(file string, params map[string]interface{}, channel chan task.TaskFeedBack) {
	var err error
	var output []string
	defer task.SendTaskOver(channel, err, output)
	ouputFormat := params["output_format"].(string)
	inputFormatStr := params["input_format"].(string)
	var inputFormat []string
	if err = json.Unmarshal([]byte(inputFormatStr), &inputFormat); err != nil {
		return
	}

	var files []string
	files, err = getAllFileOfExtensionRecursive(file, inputFormat)
	if err != nil {
		channel <- task.TaskFeedBack{Event: task.ErrorFeedBack, Message: err.Error()}
		return
	}

	var data []byte
	for _, file := range files {
		task.SendUpdate(channel, "Start converting "+file)
		newName := remplaceFileExtension(file, ouputFormat)
		cmd := exec.Command("ffmpeg", "-y", "-i", file, "-vcodec", "libx264", "-preset", "fast", "-acodec", "aac", newName, "-hider_banner")
		if data, err = cmd.CombinedOutput(); err != nil {
			return
		}
		output = append(output, newName)
		task.SendUpdate(channel, string(data))
	}

}

func remplaceFileExtension(file string, ext string) string {
	originExt := filepath.Ext(file)

	if originExt != "" {
		d := strings.TrimSuffix(file, originExt)
		println(d + ext)
		return d + ext
	}
	return file
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
