package pipeline

import (
	"errors"
	"os"

	"github.com/mitchellh/mapstructure"
)

// RPCPipeline ...
type RPCPipeline struct{}

// GetPipelines ...
func (r *RPCPipeline) GetPipelines() map[string]Pipeline {
	return Pipelines
}

// GetPipeline ...
func (r *RPCPipeline) GetPipeline(id string) Pipeline {
	return Pipelines[id]
}

// GetActive ....
func (r *RPCPipeline) GetActive(id string) *ActivePipelineStatus {
	return ActivePipelines[id]
}

// GetActives ...
func (r *RPCPipeline) GetActives() map[string]*ActivePipelineStatus {
	return ActivePipelines
}

// StartOnLocalFile ....
func (r *RPCPipeline) StartOnLocalFile(filepath string, pipelineid string, data map[string]interface{}) (status *ActivePipelineStatus) {
	var err error
	if status, err = StartOnLocalFile(filepath, pipelineid, data); err == nil {
		return status
	}
	panic(err)
}

// Create ...
func (r *RPCPipeline) Create(data map[string]interface{}) Pipeline {
	var pipeline Pipeline
	if err := mapstructure.Decode(data, &pipeline); err != nil {
		panic(err)
	}
	if _, ok := Pipelines[pipeline.ID]; ok {
		panic(errors.New("Pipeline already exists"))
	}
	if err := savePipelineFile(&pipeline); err != nil {
		panic(err)
	}
	Pipelines[pipeline.ID] = pipeline
	return pipeline
}

// Delete ...
func (r *RPCPipeline) Delete(id string) string {
	filepath := getPipelineFilePath(id)
	if err := os.Remove(filepath); err != nil {
		panic(err)
	}
	delete(Pipelines, id)
	return "OK"
}

// DeleteActive ...
func (r *RPCPipeline) DeleteActive(id string) string {
	if status, ok := ActivePipelines[id]; ok {
		if status.State != PipelineRunning {
			delete(ActivePipelines, id)
			return "OK"
		}
		panic("Pipeline is running dude")
	}
	panic("Active pipeline dont exists")
}
