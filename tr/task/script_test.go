package task

import (
	"io/ioutil"
	"testing"
)

func TestScriptTask(t *testing.T) {
	task := &InterpretorTask{
		Interpretor: "bash",
		File:        "./test.sh",
		Info:        TaskInfo{},
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

func TestScriptOp(t *testing.T) {
	data, err := ioutil.ReadFile("test.sh")
	if err != nil {
		t.Fatal(err)
	}
	ta := &InterpretorTask{
		Interpretor: "bash",
		File:        "test.sh",
		Info: TaskInfo{
			Name:        "Test script",
			Description: "Do crazy shit for nothing",
			Params: []Params{
				Params{
					Name:        "MEDIA",
					Type:        "string",
					Optional:    true,
					Description: "The subfolder to print",
				},
			},
			Return: []Return{
				Return{
					Type:        "original",
					Description: "original file",
				},
				Return{
					Type:        "file",
					Description: "a non existing file",
				},
			},
		},
	}

	err = SaveTaskScript(ta, data)
	if err != nil {
		t.Fatal(err)
	}

	if tt := GetTask("test.sh"); tt != nil {
		channel := make(chan TaskFeedBack, 5)
		tt.Execute("/home/wq/test.py", map[string]interface{}{}, channel)

		d := <-channel
		taskOver := d.Message.(TaskOver)
		for _, v := range taskOver.Files {
			println(v)
		}
	}
}
