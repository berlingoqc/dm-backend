package tasks

import "github.com/berlingoqc/dm-backend/tr"

func init() {
	tr.RegisterTask(&CPTask{})
	tr.RegisterTask(&ZipTask{})
}
