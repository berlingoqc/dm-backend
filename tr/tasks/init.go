package tasks

import "github.com/berlingoqc/dm/tr"

func init() {
	tr.RegisterTask(&CPTask{})
	tr.RegisterTask(&ZipTask{})
}
