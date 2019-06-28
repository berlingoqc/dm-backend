package pipeline

import (
	"errors"
	"io/ioutil"
	"path"
	"path/filepath"
	"strings"

	"github.com/berlingoqc/dm-backend/file"
	"github.com/berlingoqc/dm-backend/tr/task"
	"github.com/mitchellh/go-homedir"
)

// PipelineState ...
type PipelineState string

const (
	// PipelineRunning ...
	PipelineRunning PipelineState = "running"
	// PipelineCancel ...
	PipelineCancel PipelineState = "cancel"
	// PipelinePaused ...
	PipelinePaused PipelineState = "paused"
	// PipelineOver ...
	PipelineOver PipelineState = "over"
)

const (
	// OnPipelineStart ...
	OnPipelineStart = "onPipelineStart"
	// OnPipelineEnd ...
	OnPipelineEnd = "onPipelineEnd"
	// OnPipelineError ...
	OnPipelineError = "onPipelineError"
	// OnTaskStart ...
	OnTaskStart = "onTaskStart"
	// OnTaskEnd ...
	OnTaskEnd = "onTaskEnd"
	// OnTaskError ...
	OnTaskError = "onTaskError"
)

var (
	// MaximalPipelineRunning is the number of pipeline that can run at the same time
	MaximalPipelineRunning = 3
)

// ActivePipelineStatus ...
type ActivePipelineStatus struct {
	Pipeline   string                 `json:"pipeline"`
	State      PipelineState          `json:"state"`
	ActiveTask []string               `json:"activetask"`
	TaskResult map[string]interface{} `json:"tasresult"`
}

// Pipeline is a definition of task to execute on a file
type Pipeline struct {
	ID   string        `json:"id"`
	Name string        `json:"name"`
	Node task.TaskNode `json:"node"`
}

// Pipelines contains all the available pipeline
var Pipelines = make(map[string]Pipeline)

// RegisterPipeline contains the pipeline that are register
// and waiting a download before to be started
var RegisterPipeline = make(map[string]string)

// ActivePipeline contains the pipeline that are currently running
var ActivePipeline = make(map[string]*ActivePipelineStatus)

func getWorkingPath() string {
	dir, _ := homedir.Dir()
	return filepath.Join(dir, ".dm", "pipeline")
}

func getPipelineFilePath(id string) string {
	return filepath.Join(getWorkingPath(), id+".json")
}

func getPipelineFile(id string) (*Pipeline, error) {
	filepath := getPipelineFilePath(id)
	pipeline := &Pipeline{}
	return pipeline, file.LoadJSON(filepath, pipeline)
}

func savePipelineFile(pipeline *Pipeline) error {
	filepath := getPipelineFilePath(pipeline.ID)
	return file.SaveJSON(filepath, pipeline)
}

// remove from RegisterPipeline and add to ActivePipeline
func registerToActivePipeline(idRegister string, pipelineName string) *ActivePipelineStatus {
	// Delete from register pipeline
	delete(RegisterPipeline, idRegister)
	// Cree la nouvelle pipeline active
	status := &ActivePipelineStatus{Pipeline: pipelineName, State: PipelineRunning, TaskResult: make(map[string]interface{})}
	ActivePipeline[pipelineName] = status
	return status
}

func startPipeline(id string, pipeline *Pipeline) {
	status := registerToActivePipeline(id, pipeline.Name)

	currentNode := pipeline.Node
	chTasks := make(chan task.TaskFeedBack)
LoopNode:
	for {
		println("STARTING TASK ID ", currentNode.TaskID)
		task := task.GetTask(currentNode.TaskID)
		if task == nil {
			println("ERROR TASK NON TROUVER")
			break LoopNode
		}
		status.ActiveTask = append(status.ActiveTask, task.GetID())
		go task.Execute(id, currentNode.Params, chTasks)
	LoopTask:
		for {
			select {
			case feedback := <-chTasks:
				switch feedback.Event {
				case "DONE":
					println("Task done passing to next")
					if len(currentNode.NextNode) == 0 {
						break LoopNode
					} else {
						currentNode = currentNode.NextNode[0]
					}
					break LoopTask
				case "ERROR":
					println("ERROR RUNNING TASK ", currentNode.TaskID, " error : ", feedback.Message.(error).Error())
					break LoopNode
				}
			}
		}
	}
	println("FIN DE LA PIPEPINE ", id)
}

// Start ...
func Start(id string) error {
	if pipelineName, ok := RegisterPipeline[id]; ok {
		if pipeline, ok := Pipelines[pipelineName]; ok {
			startPipeline(id, &pipeline)
		}
		return errors.New("Pipeline not found " + pipelineName)
	}
	return errors.New("RegisterPipeline not found " + id)
}

func init() {
	// Loading pipeline
	folderPath := getWorkingPath()
	files, err := ioutil.ReadDir(folderPath)
	if err != nil {
		panic(err)
	}
	for _, f := range files {
		println("Loading pipeline ", f.Name())
		id := strings.TrimSuffix(f.Name(), path.Ext(f.Name()))
		pipeline, err := getPipelineFile(id)
		if err != nil {
			panic(err)
		}
		Pipelines[id] = *pipeline
	}
}
