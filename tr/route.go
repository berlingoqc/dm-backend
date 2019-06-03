package tr

// RPCTask ...
type RPCTask struct{}

// GetTasks ...
func (r *RPCTask) GetTasks() []TaskInfo {
	var ti []TaskInfo
	for _, v := range tasks {
		ti = append(ti, v.GetInfo())
	}
	return ti
}

// RPCPipeline ...
type RPCPipeline struct{}

// GetPipelines ...
func (r *RPCPipeline) GetPipelines() map[string]Pipeline {
	return Pipelines
}

// RegisterPipeline ...
func (r *RPCPipeline) RegisterPipeline() map[string]string {
	return RegisterPipeline
}
