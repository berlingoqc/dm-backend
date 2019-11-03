package tr_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/berlingoqc/dm-backend/rpcproxy"
	"github.com/berlingoqc/dm-backend/tr"
	"github.com/berlingoqc/dm-backend/tr/pipeline"
	"github.com/berlingoqc/dm-backend/tr/task"
	"github.com/berlingoqc/dm-backend/tr/triggers"

	// for the task
	_ "github.com/berlingoqc/dm-backend/tr/task/impl"
)

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
		NodeID: "abcd",
		TaskID: "copy",
		Params: map[string]string{
			"destination": testFolderPath + "/${PATH}/",
		},
		NextNode: []*task.TaskNode{},
	},
}

var pipelineData = map[string]string{
	"PATH": "figaro",
}

var channelClosing chan interface{}

var eventPipelineChannel = make(chan rpcproxy.WSMessage)

var rpc = triggers.RPC{}
var err error

func TestPipeline(t *testing.T) {
	defer func() {
		tr.StopPipelineModule()
		if err != nil {
			t.Error(err)
		}
	}()

	prepareEven()

	initModule(t)
	createPipeline(t)
	getPipeline(t)

	executePipeline(t)
}

func prepareEven() {
	pipeline.GetWorkingPath = func() string {
		return testFolderPath
	}

	_ = os.RemoveAll(testFolderPath)
	_ = os.Mkdir(testFolderPath, 0755)
}

func initModule(t *testing.T) {
	tr.InitPipelineModule(settings)

	pipeline.FeedBack = func(namespace string, event string, data interface{}) {
		fmt.Printf("NAMESPACE %s EVENT %s \n", namespace, event)
		eventPipelineChannel <- rpcproxy.WSMessage{
			Namespace: namespace,
			Data: rpcproxy.RPCCall{
				Method: event,
				Result: data,
			},
		}
	}
}

func createPipeline(t *testing.T) {
	if err = pipeline.SavePipelineFile(testPipeline); err != nil {
		panic(err)
	}

}

func getPipeline(t *testing.T) {
	if len(pipeline.Pipelines) != 1 {
		t.Fatal("No pipelines")
	}
}

func executePipeline(t *testing.T) {
	_, err := triggers.GetTrigger("manual").AddWatch("", "tr.go", &triggers.Settings{
		PipelineID: "copyPipeline",
		Data:       pipelineData,
	})
	if err != nil {
		t.Error(err)
	}

Loop:
	for {
		select {
		case e := <-eventPipelineChannel:

			if e.Data.Method == "onPipelineError" {
				t.Error("Error executing pipeline ", e.Data.Result)
			}
			if e.Data.Method == "onPipelineEnd" {
				break Loop
			}
			break
		}
	}
}
