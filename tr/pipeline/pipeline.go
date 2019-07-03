package pipeline

import (
	"errors"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/berlingoqc/dm-backend/file"
	"github.com/berlingoqc/dm-backend/tr/task"
	"github.com/mitchellh/go-homedir"
)

// FeedBack ...
var FeedBack func(namespace string, event string, data interface{})

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
	State      PipelineState       `json:"state"`
	ActiveTask []string            `json:"activetask"`
	TaskResult map[string][]string `json:"taskresult"`
	Register   RegisterPipeline    `json:"register"`
}

// RegisterPipeline ...
type RegisterPipeline struct {
	File     string        `json:"file"`
	Pipeline string        `json:"pipeline"`
	Provider string        `json:"provider"`
	Data     []interface{} `json:"data"`
}

// Pipeline is a definition of task to execute on a file
type Pipeline struct {
	ID   string        `json:"id"`
	Name string        `json:"name"`
	Node task.TaskNode `json:"node"`
}

// Pipelines contains all the available pipeline
var Pipelines = make(map[string]Pipeline)

// RegisterPipelines contains the pipeline that are register
// and waiting a download before to be started
var RegisterPipelines = make(map[string]RegisterPipeline)

// ActivePipelines contains the pipeline that are currently running
var ActivePipelines = make(map[string]*ActivePipelineStatus)

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

func createActivePipeline(id string, pipelineid string) *ActivePipelineStatus {
	status := &ActivePipelineStatus{Pipeline: pipelineid, State: PipelineRunning, TaskResult: make(map[string][]string)}
	ActivePipelines[id] = status
	FeedBack("pipeline", OnPipelineActiveUpdate, nil)
	FeedBack("pipeline", OnPipelineRegisterUpdate, nil)
	return status
}


type taskQueue struct {
	// level on the task in the pipeline
	Level int
	// Index on the task in the return
	Index int
	// The id of the task
	TaskID string
}


func startPipeline(id string, status *ActivePipelineStatus, pipeline *Pipeline) {
	var t ITask
	files := make([][]string,1)
	files[0] := []string{id}

	nextNodes := make(chan taskQueue)
	chTasks := make(chan task.TaskFeedBack)

	currentLevel := 0

	FeedBack("pipeline", OnPipelineStart, pipeline.ID)

	nextNodes <- taskQueue{
		Index: 0,
		Level: currentLevel,
		TaskID: pipeline.Node.TaskID,
	}


LoopNode:
	for {
		// Regarde dans le channel pour avoir le nom de la prochaine task a executer
		var taskQueueCurrent taskQueue
		select {
		case taskQueueCurrent = <-nextNodes:
			t = task.GetTask(taskQueueCurrent.TaskID)
			if t == nil {
				FeedBack("pipeline", OnPipelineError, "task not found")
				break LoopNode
			}
			break
		default:
			println("No more data in the channel pipeline is over")
			break LoopNode
		}
		if taskQueueCurrent.Level > currentLevel {
			println("Switching to next level of task")
			currentLevel = taskQueueCurrent.Level
		}
		// Regarde si la tache a des sous-tache et les ajoutes a la liste de tache a executer
		for _, v := range t.NextNode {

		}


		FeedBack("pipeline", OnTaskStart, currentNode.TaskID)
		status.ActiveTask = append(status.ActiveTask, t.GetID())
		go t.Execute(id, currentNode.Params, chTasks)
	LoopTask:
		for {
			select {
			case feedback := <-chTasks:
				switch feedback.Event {
				case task.DoneFeedBack:
					FeedBack("pipeline", OnTaskEnd, currentNode.TaskID)
					if len(currentNode.NextNode) == 0 {
						break LoopNode
					} else {
						if over, ok := feedback.Message.(task.TaskOver); ok && len(over.Files) > 0 {
							id = over.Files[0]
							currentNode = currentNode.NextNode[0]
						} else {
							FeedBack("pipeline", OnPipelineError, "could not get information to start next task")
						}
					}
					break LoopTask
				case task.OutFeedBack:
					FeedBack("pipeline", OnTaskUpdate, feedback.Message)
					status.TaskResult[t.GetID()] = append(status.TaskResult[t.GetID()], feedback.Message.(string))
				case task.ErrorFeedBack:
					err := feedback.Message.(error).Error()
					FeedBack("pipeline", OnTaskError, err)
					status.TaskResult[t.GetID()] = append(status.TaskResult[t.GetID()], err)
					break LoopNode
				}
			}
		}
	}
	status.State = PipelineOver
	FeedBack("pipeline", OnPipelineEnd, id)

}

// StartOnLocalFile ...
func StartOnLocalFile(filepath string, pipelineid string) (*ActivePipelineStatus, error) {
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return nil, err
	}
	if pipeline, ok := Pipelines[pipelineid]; ok {
		status := &ActivePipelineStatus{Pipeline: pipeline.ID, State: PipelineRunning, TaskResult: make(map[string][]string)}
		ActivePipelines[filepath] = status
		return status, nil
	}
	return nil, errors.New("Pipeline not found")
}

// StartFromRegister ...
func StartFromRegister(id string) (*ActivePipelineStatus, error) {
	if pipelineName, ok := RegisterPipelines[id]; ok {
		if pipeline, ok := Pipelines[pipelineName.Pipeline]; ok {
			delete(RegisterPipelines, id)
			status := createActivePipeline(id, pipeline.ID)
			startPipeline(id, status, &pipeline)
			return status, nil
		}
		return nil, errors.New("Pipeline not found " + pipelineName.Pipeline)
	}
	return nil, errors.New("RegisterPipeline not found " + id)
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
