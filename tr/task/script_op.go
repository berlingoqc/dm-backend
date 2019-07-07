package task

import (
	"errors"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/berlingoqc/dm-backend/file"
)

// SaveTaskScript ...
func SaveTaskScript(task *InterpretorTask, fileContent []byte) error {
	f := filepath.Join(getWorkingPath(), task.Interpretor)
	if _, err := os.Stat(f); os.IsNotExist(err) {
		if err := os.Mkdir(f, 0755); err != nil {
			return err
		}
	}
	if err := file.SaveJSON(filepath.Join(f, task.GetID()+".json"), task.GetInfo()); err != nil {
		return err
	}
	if err := ioutil.WriteFile(filepath.Join(f, task.File), fileContent, 0644); err != nil {
		return err
	}
	tasks[task.GetID()] = task
	return nil
}

// DeleteTaskScript ...
func DeleteTaskScript(taskid string) error {
	if d, ok := tasks[taskid].(*InterpretorTask); ok {

		folder := filepath.Join(getWorkingPath(), d.Interpretor)
		if err := os.Remove(filepath.Join(folder, d.GetID())); err != nil {
			return err
		}
		if err := os.Remove(filepath.Join(folder, d.GetID()+".json")); err != nil {
			return err
		}
		delete(tasks, taskid)
		return nil
	}
	return errors.New("")
}

// GetAllTaskScript ...
func GetAllTaskScript() (ret []*InterpretorTask, err error) {
	scriptFolder := getWorkingPath()
	if files, err := ioutil.ReadDir(scriptFolder); err == nil {
		for _, f := range files {
			if f.IsDir() {
				interpretor := f.Name()
				folder := path.Join(scriptFolder, interpretor, "*.json")
				if filesConfig, err := filepath.Glob(folder); err == nil {
					info := &TaskInfo{}
					for _, ff := range filesConfig {
						if err = file.LoadJSON(ff, info); err != nil {
							return nil, err
						}
						task := &InterpretorTask{
							Interpretor: interpretor,
							File:        info.Name,
							Info:        *info,
						}
						ret = append(ret, task)
					}
				}
			}
		}
	} else {
		return nil, err
	}
	return ret, nil
}
