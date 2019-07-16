package tr

import "github.com/fsnotify/fsnotify"

var watcher *fsnotify.Watcher

func CreateFileWatcher() {
	var err error
	watcher, err = fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}
	defer watcher.Close()

	go func() {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			println("event: ",event)
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			println("error: ",err.Error())
		}
	}()
}

func AddFileWatch(file string) error {
	if watcher == nil {
		CreateFileWatcher()
	}
	return watcher.Add(file)
}

func RemoveFileWatch(file string) {
	return watcher.Remove(file)
}

func CloseFileWatch() {
	return watcher.Close()
}