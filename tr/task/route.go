package task

import "github.com/mitchellh/mapstructure"

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

// SaveTaskScript ...
func (r *RPCTask) SaveTaskScript(taskIn map[string]interface{}, data []byte) {
	task := &InterpretorTask{}
	if err := mapstructure.Decode(data, task); err != nil {
		panic(err)
	}
	if err := SaveTaskScript(task, data); err != nil {
		panic(err)
	}
}

// DeleteTaskScript ...
func (r *RPCTask) DeleteTaskScript(taskid string) {
	if err := DeleteTaskScript(taskid); err != nil {
		panic(err)
	}
}
