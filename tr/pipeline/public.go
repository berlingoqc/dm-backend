package pipeline

import (
	"errors"
	"os"

	"github.com/berlingoqc/dm-backend/tr/task"
)

// FeedBack ...
var FeedBack func(namespace string, event string, data interface{})

// State ...
type State string

const (
	// PipelineRunning ...
	PipelineRunning State = "running"
	// PipelineCancel ...
	PipelineCancel State = "cancel"
	// PipelinePaused ...
	PipelinePaused State = "paused"
	// PipelineOver ...
	PipelineOver State = "over"
)

const (
	// OnPipelineStart ...
	OnPipelineStart = "onPipelineStart"
	// OnPipelineEnd ...
	OnPipelineEnd = "onPipelineEnd"
	// OnPipelineError ...
	OnPipelineError = "onPipelineError"
	// OnPipelineStatusUpdate ...
	OnPipelineStatusUpdate = "onPipelineStatusUpdate"
	// OnTaskStart ...
	OnTaskStart = "onTaskStart"
	// OnTaskEnd ...
	OnTaskEnd = "onTaskEnd"
	// OnTaskError ...
	OnTaskError = "onTaskError"
	// OnTaskUpdate ...
	OnTaskUpdate = "onTaskUpdate"
	// OnPipelineActiveUpdate ...
	OnPipelineActiveUpdate = "onPipelineActiveUpdate"
)

func eventOnPipelineStart(pipeline *ActivePipelineStatus) {
	FeedBack("pipeline", OnPipelineStart, pipeline)
}

func eventOnPipelineEnd(status *ActivePipelineStatus) {
	FeedBack("pipeline", OnPipelineEnd, status)
}

func eventOnPipelineError(err error) {
	FeedBack("pipeline", OnPipelineError, err.Error())
}

func eventOnPipelineStatusUpdate(status *ActivePipelineStatus) {
	FeedBack("pipeline", OnPipelineActiveUpdate, status)
}

func eventOnTaskStart(id string) {
	FeedBack("pipeline", OnTaskStart, id)
}

func eventOnTaskEnd(id string, files []string) {
	FeedBack("pipeline", OnTaskEnd, struct {
		id    string
		files []string
	}{id, files})
}

func eventOnTaskError(id string, err error, output []string) {
	FeedBack("pipeline", OnTaskError, struct {
		id     string
		err    string
		output []string
	}{id, err.Error(), output})

}

func eventOnTaskUpdate(id string, output []string) {
	FeedBack("pipeline", OnTaskUpdate, struct {
		id     string
		output []string
	}{id, output})
}

func eventOnPipelineActiveUpdate() {
	FeedBack("pipeline", OnPipelineActiveUpdate, ActivePipelines)
}

var (
	// MaximalPipelineRunning is the number of pipeline that can run at the same time
	MaximalPipelineRunning = 3
)

// ActivePipelineStatus ...
type ActivePipelineStatus struct {
	File       string              `json:"file"`
	Pipeline   string              `json:"pipeline"`
	State      State               `json:"state"`
	ActiveTask string              `json:"activetask"`
	TaskOutput map[string][]string `json:"taskouput"`
	TaskResult map[string][]string `json:"taskresult"`
}

// Variables ...
type Variables struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
}

// Pipeline is a definition of task to execute on a file
type Pipeline struct {
	ID        string         `json:"id"`
	Name      string         `json:"name"`
	Variables []Variables    `json:"variables"`
	Node      *task.TaskNode `json:"node"`
}

// Pipelines contains all the available pipeline
var Pipelines = make(map[string]Pipeline)

// ActivePipelines contains the pipeline that are currently running
var ActivePipelines = make(map[string]*ActivePipelineStatus)

func StartOnLocalFile(filepath string, pipelineid string, data map[string]interface{}) (*ActivePipelineStatus, error) {
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return nil, err
	}
	if pipeline, ok := Pipelines[pipelineid]; ok {
		return startPipeline(filepath, &pipeline, data)
	}
	return nil, errors.New("Pipeline not found")
}
