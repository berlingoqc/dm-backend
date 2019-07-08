package pipeline

import (
	"path/filepath"

	"github.com/berlingoqc/dm-backend/file"
	"github.com/berlingoqc/dm-backend/tr/task"
	"github.com/mitchellh/go-homedir"
)

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
	status := &ActivePipelineStatus{Pipeline: pipelineid, State: PipelineRunning, File: id, TaskResult: make(map[string][]string), TaskOutput: make(map[string][]string)}
	ActivePipelines[id] = status
	FeedBack("pipeline", OnPipelineActiveUpdate, nil)
	FeedBack("pipeline", OnPipelineRegisterUpdate, id)
	return status
}

type taskQueue struct {
	// Index on the task in the return
	Index int `json:"index"`
	// The task node
	TaskNode *task.TaskNode `json:"tasknode"`
	// Result of the task , the list of file
	Result []string `json:"result"`
	// Previous item to retrieve the result
	Previous *taskQueue `json:"previous"`
}

func startPipeline(id string, pip *Pipeline, data map[string]interface{}) (*ActivePipelineStatus, error) {
	newPipeline := &Pipeline{}
	cloneValue(pip, newPipeline)

	status := createActivePipeline(id, pip.ID)

	go pipeline(id, status, pip, data)
	return status, nil
}

func pipeline(id string, status *ActivePipelineStatus, pipeline *Pipeline, data map[string]interface{}) {
	var t task.ITask

	nextNodes := make(chan *taskQueue, 5)
	chTasks := make(chan task.TaskFeedBack)
	FeedBack("pipeline", OnPipelineStart, pipeline.ID)

	nextNodes <- &taskQueue{
		Index:    0,
		TaskNode: pipeline.Node,
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
		status.ActiveTask = t.GetID()
		// Get le nom du next fichier
		if taskQueueCurrent.Previous == nil || len(taskQueueCurrent.Previous.Result) < taskQueueCurrent.Index {
			FeedBack("pipeline", OnPipelineError, "files needed at index is not present")
			break LoopNode
		}
		nextFile := taskQueueCurrent.Previous.Result[taskQueueCurrent.Index]
		println("Starting task on file ", nextFile)

		params, _ := replaceParams(currentNode.Params, data)

		go t.Execute(nextFile, params, chTasks)

	LoopTask:
		for {
			select {
			case feedback := <-chTasks:
				switch feedback.Event {
				case task.DoneFeedBack:
					if over, ok := feedback.Message.(task.TaskOver); ok && len(over.Files) > 0 {
						taskQueueCurrent.Result = over.Files
						status.TaskOutput[currentNode.TaskID] = over.Files
						status.ActiveTask = ""
						FeedBack("pipeline", OnTaskEnd, status)
						break LoopTask
					} else {
						FeedBack("pipeline", OnPipelineError, "could not get information to start next task")
						break LoopNode
					}
				case task.OutFeedBack:
					FeedBack("pipeline", OnTaskUpdate, feedback.Message)
					status.TaskResult[currentNode.TaskID] = append(status.TaskResult[t.GetID()], feedback.Message.(string))
				case task.ErrorFeedBack:
					err := feedback.Message.(error).Error()
					FeedBack("pipeline", OnTaskError, err)
					status.TaskResult[currentNode.TaskID] = append(status.TaskResult[t.GetID()], err)
					break LoopNode
				}
			}
		}
	}
	status.State = PipelineOver
	FeedBack("pipeline", OnPipelineEnd, id)
}
