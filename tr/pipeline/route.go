package pipeline

import (
	"errors"
	"os"

	"github.com/mitchellh/mapstructure"
)

// MapValueToSlice ...
func MapValueToSlice(data map[string]interface{}, d []interface{}) {
	for _, v := range data {
		d = append(d, v)
	}
}

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

// GetRegister ...
func (r *RPCPipeline) GetRegister() map[string]RegisterPipeline {
	return RegisterPipelines
}

// GetActive ...
func (r *RPCPipeline) GetActive() []*ActivePipelineStatus {
	var a []*ActivePipelineStatus
	for _, v := range ActivePipelines {
		a = append(a, v)
	}
	return a
}

// Register ...
func (r *RPCPipeline) Register(handlerName, pipeline string, data []interface{}) {
	if handler, ok := Handlers[handlerName]; ok {
		filepath, err := handler.GetFilePath(data)
		if err != nil {
			panic(err)
		}
		if filepath == "" {
			panic("Filepath from handler " + handlerName + " is empty")
		}
		RegisterPipelines[filepath] = RegisterPipeline{
			File:     filepath,
			Pipeline: pipeline,
			Provider: handlerName,
			Data:     data,
		}
		println("Pipeline register")
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

// DeleteRegister ...
func (r *RPCPipeline) DeleteRegister(id string) {
	delete(RegisterPipelines, id)
}
