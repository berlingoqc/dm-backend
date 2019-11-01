package triggers

import (
	"github.com/fsnotify/fsnotify"
)

var watcher *fsnotify.Watcher
var watchFiles []string

// FWEvent ...
type FWEvent string

var (
	// FWEventWrite ...
	FWEventWrite = "eventWrite"
	// FWEventUpdate ...
	FWEventUpdate = "eventUpdate"
	// FWEventCreated ...
	FWEventCreated = "eventCreated"
	// FWEventDelete ...
	FWEventDelete = "eventDelete"
)

func deleteItem(slice []string, item string) (ret []string) {
	for _, i := range slice {
		if i != item {
			ret = append(ret, i)
		}
	}
	return ret
}

// RemoveFileWatch ...
func RemoveFileWatch(file string) error {
	watchFiles = deleteItem(watchFiles, file)
	return watcher.Remove(file)
}

// FileWatchTrigger ...
type FileWatchTrigger struct {
	watcher    *fsnotify.Watcher
	filesWatch map[int64]WatchInfo
}

// AddWatch ...
func (f *FileWatchTrigger) AddWatch(event string, param interface{}, settings *Settings) (int64, error) {
	file := param.(string)
	f.filesWatch[getID()] = WatchInfo{
		Trigger:  "file_watch",
		Event:    event,
		Param:    param,
		Settings: settings,
	}
	return 0, f.watcher.Add(file)
}

// DeleteWatch ...
func (f *FileWatchTrigger) DeleteWatch(id int64) error {
	return f.watcher.Remove("")
}

// GetWatchInfo ...
func (f *FileWatchTrigger) GetWatchInfo() *map[int64]WatchInfo {
	return nil
}

// Init ...
func (f *FileWatchTrigger) Init(ch chan PipelineTrigger, signal chan interface{}) {
	var err error
	f.watcher, err = fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}
	go func() {
		defer func() {
			if err := f.watcher.Close(); err != nil {
				panic(err)
			}
		}()
		for {
			select {
			case event, ok := <-f.watcher.Events:
				if !ok {
					return
				}
				// Regarde nos events qu'on n'a et trigger celle qui faut
				for i, k := range f.filesWatch {
					if k.Param.(string) == getEvent(event.Op) {
						println("Event for pipeline")

						ch <- PipelineTrigger{
							File:       event.Name,
							PipelineID: k.Settings.PipelineID,
							Data:       k.Settings.Data,
						}

						if k.Settings.RemoveAfterRun {
							delete(f.filesWatch, i)
						}
					}
				}

				break

			case err, ok := <-f.watcher.Errors:
				if !ok {
					return
				}
				println("error file watch ", err.Error())
				break
			case _ = <-signal:
				println("Over baby blue file watch")
				return
			}
		}
	}()
}

func getEvent(op fsnotify.Op) string {
	if op&fsnotify.Create == fsnotify.Create {
		return FWEventCreated
	}
	if op&fsnotify.Chmod == fsnotify.Chmod {
		return FWEventUpdate
	}
	if op&fsnotify.Remove == fsnotify.Remove {
		return FWEventDelete
	}
	panic("EVENT NOT SPORTED")
}
