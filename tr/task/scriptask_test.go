package task

import "testing"

func TestScriptTask(t *testing.T) {
	task := &InterpretorTask{
		Interpretor: "bash",
		File:        "./test.sh",
		info:        TaskInfo{},
	}

	channel := make(chan TaskFeedBack, 5)

	task.Execute("/home/wq/test.py", map[string]interface{}{
		"MEDIA": "lol",
	}, channel)

	d := <-channel
	taskOver := d.Message.(TaskOver)
	for _, v := range taskOver.Files {
		println(v)
	}
}
