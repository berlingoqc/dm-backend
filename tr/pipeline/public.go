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
	// OnPipelineRegisterUpdate ...
	OnPipelineRegisterUpdate = "onPipelineRegisterUpdate"
	// OnPipelineActiveUpdate ...
	OnPipelineActiveUpdate = "onPipelineActiveUpdate"
)

var (
	// MaximalPipelineRunning is the number of pipeline that can run at the same time
	MaximalPipelineRunning = 3
)

// ActivePipelineStatus ...
type ActivePipelineStatus struct {
	Pipeline   string              `json:"pipeline"`
	State      State               `json:"state"`
	ActiveTask []string            `json:"activetask"`
	TaskResult map[string][]string `json:"taskresult"`
	Register   RegisterPipeline    `json:"register"`
}

// RegisterPipeline ...
type RegisterPipeline struct {
	File     string                 `json:"file"`
	Pipeline string                 `json:"pipeline"`
	Provider string                 `json:"provider"`
	Data     map[string]interface{} `json:"data"`
}

// Variables ...
type Variables struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
}

// Pipeline is a definition of task to execute on a file
type Pipeline struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Variables []Variables
	Node      *task.TaskNode `json:"node"`
}

// Pipelines contains all the available pipeline
var Pipelines = make(map[string]Pipeline)

// RegisterPipelines contains the pipeline that are register
// and waiting a download before to be started
var RegisterPipelines = make(map[string]RegisterPipeline)

// ActivePipelines contains the pipeline that are currently running
var ActivePipelines = make(map[string]*ActivePipelineStatus)

// StartOnLocalFile ...
func StartOnLocalFile(filepath string, pipelineid string, data map[string]interface{}) (*ActivePipelineStatus, error) {
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return nil, err
	}
	if pipeline, ok := Pipelines[pipelineid]; ok {
		return startPipeline(filepath, &pipeline, data)
	}
	return nil, errors.New("Pipeline not found")
}

// StartFromRegister ...
func StartFromRegister(id string) (*ActivePipelineStatus, error) {
	if pipelineName, ok := RegisterPipelines[id]; ok {
		if pipeline, ok := Pipelines[pipelineName.Pipeline]; ok {
			delete(RegisterPipelines, id)
			return startPipeline(id, &pipeline, pipelineName.Data)
		}
		return nil, errors.New("Pipeline not found " + pipelineName.Pipeline)
	}
	return nil, errors.New("RegisterPipeline not found " + id)
}