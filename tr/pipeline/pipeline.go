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
	// Index on the task in the return
	Index int
	// The task node
	TaskNode *task.TaskNode
	// Result of the task , the list of file
	Result []string
	// Previous item to retrieve the result
	Previous *taskQueue
}

func startPipeline(id string, status *ActivePipelineStatus, pipeline *Pipeline) {
	var t task.ITask

	nextNodes := make(chan *taskQueue, 5)
	chTasks := make(chan task.TaskFeedBack)
	FeedBack("pipeline", OnPipelineStart, pipeline.ID)

	nextNodes <- &taskQueue{
		Index:    0,
		TaskNode: &pipeline.Node,
		Previous: &taskQueue{
			Result: []string{id},
		},
	}

	var taskQueueCurrent *taskQueue
LoopNode:
	for {
		// Regarde dans le channel pour avoir le nom de la prochaine task a executer
		select {
		case taskQueueCurrent = <-nextNodes:
			t = task.GetTask(taskQueueCurrent.TaskNode.TaskID)
			if t == nil {
				FeedBack("pipeline", OnPipelineError, "task not found")
				break LoopNode
			}
			break
		default:
			println("No more data in the channel pipeline is over")
			break LoopNode
		}
		// Regarde si la tache a des sous-tache et les ajoutes a la liste de tache a executer
		for i, v := range taskQueueCurrent.TaskNode.NextNode {
			if v != nil {
				println("Adding task ", v.TaskID)
				nextNodes <- &taskQueue{
					Index:    i,
					TaskNode: v,
					Previous: taskQueueCurrent,
				}
			}
		}
		currentNode := taskQueueCurrent.TaskNode

		FeedBack("pipeline", OnTaskStart, currentNode.TaskID)
		status.ActiveTask = append(status.ActiveTask, t.GetID())
		// Get le nom du next fichier
		if taskQueueCurrent.Previous == nil || len(taskQueueCurrent.Previous.Result) < taskQueueCurrent.Index {
			FeedBack("pipeline", OnPipelineError, "files needed at index is not present")
			break LoopNode
		}
		nextFile := taskQueueCurrent.Previous.Result[taskQueueCurrent.Index]
		println("Starting task on file ", nextFile)

		go t.Execute(nextFile, currentNode.Params, chTasks)
	LoopTask:
		for {
			select {
			case feedback := <-chTasks:
				switch feedback.Event {
				case task.DoneFeedBack:
					FeedBack("pipeline", OnTaskEnd, currentNode.TaskID)
					if over, ok := feedback.Message.(task.TaskOver); ok && len(over.Files) > 0 {
						taskQueueCurrent.Result = over.Files
						break LoopTask
					} else {
						FeedBack("pipeline", OnPipelineError, "could not get information to start next task")
						break LoopNode
					}
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
		startPipeline(filepath, status, &pipeline)
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
