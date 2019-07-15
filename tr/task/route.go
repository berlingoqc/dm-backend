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

// InputSaveTaskScript ...
type InputSaveTaskScript struct {
	TaskIn InterpretorTask `json:"taskin"`
	Data   []byte          `json:"data"`
}

// SaveTaskScript ...
func (r *RPCTask) SaveTaskScript(taskIn map[string]interface{}) string {
	task := &InputSaveTaskScript{}
	if err := mapstructure.Decode(taskIn, task); err != nil {
		panic(err)
	}
	if err := SaveTaskScript(&task.TaskIn, task.Data); err != nil {
		panic(err)
	}
	return "OK"
}

// DeleteTaskScript ...
func (r *RPCTask) DeleteTaskScript(taskid string) string {
	if err := DeleteTaskScript(taskid); err != nil {
		panic(err)
	}
	return "OK"
}
