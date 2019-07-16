package tr

import "github.com/fsnotify/fsnotify"

var watcher *fsnotify.Watcher
var watchFiles []string

func deleteItem(slice []string, item string) (ret []string) {
	for _, i := range slice {
		if i != item {
			ret = append(ret, i)
		}
	}
	return ret
}

func CreateFileWatcher() error {
	var err error
	watcher, err = fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	go func() {
		defer watcher.Close()
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					TriggerEventChannel <- TriggerEvent{
						File:  event.Name,
						Event: "onFileWrite",
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				println("error: ", err.Error())
			}
		}
	}()

	return nil
}

func AddFileWatch(file string) error {
	if watcher == nil {
		CreateFileWatcher()
	}
	watchFiles = append(watchFiles, file)
	return watcher.Add(file)
}

func RemoveFileWatch(file string) error {
	watchFiles = deleteItem(watchFiles, file)
	return watcher.Remove(file)
}

func CloseFileWatch() error {
	return watcher.Close()
}
