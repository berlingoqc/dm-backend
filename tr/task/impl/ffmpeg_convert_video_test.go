package tasks

import (
	"testing"

	"github.com/berlingoqc/dm-backend/tr/task"
)

var folderVideo = "/var/share/Media/Movie"

func TestTaskFFMPEGConvert(t *testing.T) {
	ch := make(chan task.TaskFeedBack, 1)
	ffmpegConverter := &FFMPEGConvertVideo{}

	params := map[string]interface{}{
		"output_format": "mp4",
		"input_format":  `[".avi",".wav"]`,
	}

	ffmpegConverter.Execute(folderVideo, params, ch)

}
