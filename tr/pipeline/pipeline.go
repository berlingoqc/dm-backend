package pipeline

import (
	"errors"
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
	eventOnPipelineActiveUpdate()
	eventOnPipelineRegisterUpdate()
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
	eventOnPipelineStart(status)

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
				eventOnPipelineError(errors.New("Task not found " + taskQueueCurrent.TaskNode.TaskID))
				break LoopNode
			}
			break
		default:
			break LoopNode
		}
		// Regarde si la tache a des sous-tache et les ajoutes a la liste de tache a executer
		for i, v := range taskQueueCurrent.TaskNode.NextNode {
			if v != nil {
				nextNodes <- &taskQueue{
					Index:    i,
					TaskNode: v,
					Previous: taskQueueCurrent,
				}
			}
		}
		currentNode := taskQueueCurrent.TaskNode

		eventOnTaskStart(currentNode.NodeID)

		status.ActiveTask = t.GetID()
		// Get le nom du next fichier
		if taskQueueCurrent.Previous == nil || len(taskQueueCurrent.Previous.Result) < taskQueueCurrent.Index {
			eventOnPipelineError(errors.New("files needed at index is not present"))
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
						eventOnTaskEnd(currentNode.NodeID, taskQueueCurrent.Result)
						break LoopTask
					} else {
						eventOnPipelineError(errors.New("could not get information to start next task from output"))
						break LoopNode
					}
				case task.OutFeedBack:
					status.TaskResult[currentNode.TaskID] = append(status.TaskResult[t.GetID()], feedback.Message.(string))
					eventOnTaskUpdate(currentNode.NodeID, status.TaskResult[currentNode.TaskID])
				case task.ErrorFeedBack:
					err := feedback.Message.(error)
					status.TaskResult[currentNode.TaskID] = append(status.TaskResult[t.GetID()], err.Error())
					eventOnTaskError(currentNode.NodeID, err, status.TaskResult[currentNode.TaskID])
					break LoopNode
				}
			}
		}
	}
	status.State = PipelineOver
	eventOnPipelineEnd(status)
}
