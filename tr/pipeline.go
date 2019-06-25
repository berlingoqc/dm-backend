package tr

import (
	"io/ioutil"
	"path"
	"path/filepath"
	"strings"

	"github.com/berlingoqc/dm-backend/file"
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

var (
	// MaximalPipelineRunning is the number of pipeline that can run at the same time
	MaximalPipelineRunning = 3
)

// ActivePipelineStatus ...
type ActivePipelineStatus struct {
	Pipeline   string
	State      PipelineState
	ActiveTask []string
	TaskResult map[string]interface{}
}

// Pipeline is a definition of task to execute on a file
type Pipeline struct {
	ID   string
	Name string
	Node TaskNode
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

func startPipeline(id string) {
	if pipelineName, ok := RegisterPipeline[id]; ok {
		if pipeline, ok := Pipelines[pipelineName]; ok {
			delete(RegisterPipeline, id)
			status := &ActivePipelineStatus{Pipeline: pipelineName, State: PipelineRunning, TaskResult: make(map[string]interface{})}
			ActivePipeline[pipelineName] = status
			currentNode := pipeline.Node
			chTasks := make(chan TaskFeedBack)
		LoopNode:
			for {
				println("STARTING TASK ID ", currentNode.TaskID)
				task := GetTask(currentNode.TaskID)
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
		} else {
			println("PIPELINE NOT FOUND")
		}
	} else {
		println("REGISTER PIPELINE NOT FOUND ", id)
	}

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
