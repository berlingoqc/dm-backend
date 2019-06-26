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

// GetRegister ...
func (r *RPCPipeline) GetRegister() map[string]string {
	return RegisterPipeline
}

// GetActive ...
func (r *RPCPipeline) GetActive() map[string]*ActivePipelineStatus {
	return ActivePipeline
}

// Register ...
func (r *RPCPipeline) Register(handler, pipeline string, data []interface{}) {
	if handler, ok := Handlers[handler]; ok {
		filepath, err := handler.GetFilePath(data)
		if err != nil {
			panic(err)
		}
		RegisterPipeline[filepath] = pipeline
	} else {
		panic(errors.New("Cant find handler"))
	}
}

// Create ...
func (r *RPCPipeline) Create(data map[string]interface{}) {
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
}

// Delete ...
func (r *RPCPipeline) Delete(id string) {
	filepath := getPipelineFilePath(id)
	if err := os.Remove(filepath); err != nil {
		panic(err)
	}
	delete(Pipelines, id)
}
