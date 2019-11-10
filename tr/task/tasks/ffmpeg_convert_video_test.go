package tasks_test

import (
	"testing"

	"github.com/berlingoqc/dm-backend/tr/task"
	"github.com/berlingoqc/dm-backend/tr/task/tasks"
)

var folderVideo = "/var/share/Media/Show/The Office US/Season.1"

func TestTaskFFMPEGConvert(t *testing.T) {
	ch := make(chan task.TaskFeedBack, 1)
	ffmpegConverter := &tasks.FFMPEGConvertVideo{}

	params := map[string]interface{}{
		"output_format": ".mp4",
		"input_format":  `[".avi",".wav",".mkv"]`,
	}

	ffmpegConverter.Execute(folderVideo, params, ch)

}
