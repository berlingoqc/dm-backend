package triggers

// RPC ...
type RPC struct{}

// TriggerRegister to manually trigger a register task
func (t *RPC) TriggerRegister(event, file string) {
	TriggerEventChannel <- TriggerEvent{
		Event: event,
		File:  file,
	}
}

// FileWatchRPC ...
type FileWatchRPC struct{}

func (f *FileWatchRPC) AddFile(file string) string {
	if err := AddFileWatch(file); err != nil {
		panic(err)
	}
	return "OK"
}

func (f *FileWatchRPC) RemoveFile(file string) string {
	if err := RemoveFileWatch(file); err != nil {
		panic(err)
	}
	return "OK"
}

func (f *FileWatchRPC) Files() []string {
	return watchFiles
}
