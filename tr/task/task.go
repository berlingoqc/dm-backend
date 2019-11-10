package task

// Params ...
type Params struct {
	Name         string `json:"name"`
	Type         string `json:"type"`
	Optional     bool   `json:"optional"`
	DefaultValue string `json:"default_value"`
	Description  string `json:"description"`
}

// Return ...
type Return struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}

// TaskInfo ...
type TaskInfo struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Provider    string   `json:"provider"`
	Params      []Params `json:"params"`
	Return      []Return `json:"return"`
}

// TaskOver ...
type TaskOver struct {
	Files []string
}

// TypeTaskFeedBack ...
type TypeTaskFeedBack string

const (
	// ErrorFeedBack ...
	ErrorFeedBack TypeTaskFeedBack = "ERROR"
	// DoneFeedBack ...
	DoneFeedBack TypeTaskFeedBack = "DONE"
	// OutFeedBack ...
	OutFeedBack TypeTaskFeedBack = "OUTPUT"
)

// TaskFeedBack ...
type TaskFeedBack struct {
	Event   TypeTaskFeedBack
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
	NodeID   string            `json:"nodeid"`
	TaskID   string            `json:"taskid"`
	Params   map[string]string `json:"params"`
	NextNode []*TaskNode       `json:"nextnode"`
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
func SendDone(ch chan TaskFeedBack, returnFiles []string) {
	ch <- TaskFeedBack{
		Event: "DONE",
		Message: TaskOver{
			Files: returnFiles,
		},
	}
}

// SendUpdate ...
func SendUpdate(ch chan TaskFeedBack, msg interface{}) {
	ch <- TaskFeedBack{
		Event:   OutFeedBack,
		Message: msg,
	}
}

// SendTaskOver ...
func SendTaskOver(ch chan TaskFeedBack, err error, files []string) {
	if err != nil {
		SendError(ch, err)
	} else {
		SendDone(ch, files)
	}
}
