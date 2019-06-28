package task

// Params ...
type Params struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Optional    bool   `json:"optional"`
	Description string `json:"description"`
}

// TaskInfo ...
type TaskInfo struct {
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	Params       []Params `json:"params"`
	NumberReturn int      `json:"numberreturn"`
}

// TaskOver ...
type TaskOver struct {
	Files []string
}

// TaskFeedBack ...
type TaskFeedBack struct {
	Event   string
	Message interface{}
}

// ITask  ...
type ITask interface {
	Get() ITask
	GetID() string
	GetInfo() TaskInfo
	Execute(file string, params map[string]interface{}, channel chan TaskFeedBack)
}

// TaskNode ...
type TaskNode struct {
	TaskID   string                 `json:"taskid"`
	Params   map[string]interface{} `json:"params"`
	NextNode []TaskNode             `json:"nextnode"`
}

var tasks = make(map[string]ITask)

// GetTask ...
func GetTask(task string) ITask {
	if t, ok := tasks[task]; ok {
		return t.Get()
	}
	return nil
}

// RegisterTask ...
func RegisterTask(task ITask) {
	println("Registering task ", task.GetID)
	tasks[task.GetInfo().Name] = task
}

// SendError ...
func SendError(ch chan TaskFeedBack, er error) {
	if er != nil {
		ch <- TaskFeedBack{
			Event:   "ERROR",
			Message: er,
		}
	}
}

// SendDone ...
func SendDone(ch chan TaskFeedBack, msg interface{}) {
	ch <- TaskFeedBack{
		Event:   "DONE",
		Message: msg,
	}
}
