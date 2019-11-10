package tasks

import "github.com/berlingoqc/dm-backend/tr/task"

func init() {
	task.RegisterTask(&CPTask{})
	task.RegisterTask(&ZipTask{})
	task.RegisterTask(&SleepTask{})
	task.RegisterTask(&FFMPEGConvertVideo{})

	tasks, err := task.GetAllTaskScript()
	if err != nil {
		panic(err)
	}
	for _, t := range tasks {
		task.RegisterTask(t)
	}

}
