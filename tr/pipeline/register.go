package pipeline

import "errors"

// RegisterPipeline ...
type RegisterPipeline struct {
	File       string                 `json:"file"`
	Pipeline   string                 `json:"pipeline"`
	Provider   string                 `json:"provider"`
	AutoRemove bool                   `json:"autoremove"`
	Data       map[string]interface{} `json:"data"`
}

// RegisterPipelines contains the pipeline that are register
// and waiting a download before to be started
var RegisterPipelines = make(map[string]RegisterPipeline)

// StartFromRegister ...
func StartFromRegister(id string) (*ActivePipelineStatus, error) {
	if pipelineName, ok := RegisterPipelines[id]; ok {
		if pipeline, ok := Pipelines[pipelineName.Pipeline]; ok {
			if pipelineName.AutoRemove {
				delete(RegisterPipelines, id)
			}
			return startPipeline(id, &pipeline, pipelineName.Data)
		}
		return nil, errors.New("Pipeline not found " + pipelineName.Pipeline)
	}
	return nil, errors.New("RegisterPipeline not found " + id)
}
