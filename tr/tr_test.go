package tr_test

import (
	"os"
	"testing"

	"github.com/berlingoqc/dm-backend/tr"
	"github.com/berlingoqc/dm-backend/tr/pipeline"
	"github.com/berlingoqc/dm-backend/tr/task"

	// for the task
	_ "github.com/berlingoqc/dm-backend/tr/task/impl"
)

func TestPipeline(t *testing.T) {
	defer destroyEven()

	prepareEven()

	initModule(t)
	createPipeline(t)
	getPipeline(t)
}

const testFolderPath = "./pipeline_test"

var settings = tr.Settings{
	ConcurrentPipeline: 1,
}

var testPipeline = &pipeline.Pipeline{
	ID:   "copyPipeline",
	Name: "Copy to another location twice",
	Variables: []pipeline.Variables{
		pipeline.Variables{
			Name:        "PATH",
			Type:        "string",
			Description: "Boom boom boom",
		},
	},
	Node: &task.TaskNode{
		NodeID:   "abcd",
		TaskID:   "cp",
		Params:   map[string]string{},
		NextNode: []*task.TaskNode{},
	},
}

var channelClosing chan interface{}

func prepareEven() {
	pipeline.GetWorkingPath = func() string {
		return testFolderPath
	}

	_ = os.RemoveAll(testFolderPath)
	_ = os.Mkdir(testFolderPath, 0755)
}

func destroyEven() {

	channelClosing <- 0
}

func initModule(t *testing.T) {
	tr.InitPipelineModule(settings)
}

func createPipeline(t *testing.T) {

}

func getPipeline(t *testing.T) {

}
