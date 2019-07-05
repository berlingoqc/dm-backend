package task

import "testing"

func TestScriptTask(t *testing.T) {
	task := &InterpretorTask{
		Interpretor: "bash",
		File:        "./test.sh",
		info:        TaskInfo{},
	}

	channel := make(chan TaskFeedBack, 5)

	task.Execute("/home/wq/test.py", nil, channel)

	d := <-channel
	taskOver := d.Message.(TaskOver)
	println(taskOver.Files[0])
}
