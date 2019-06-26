package task

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
