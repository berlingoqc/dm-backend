package tasks

import "github.com/berlingoqc/dm-backend/tr/task"

func init() {
	task.RegisterTask(&CPTask{})
	task.RegisterTask(&ZipTask{})
}
