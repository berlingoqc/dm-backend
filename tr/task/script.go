package task

import (
	"errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/mitchellh/go-homedir"
)

func getWorkingPath() string {
	dir, _ := homedir.Dir()
	return filepath.Join(dir, ".dm", "script")
}

func removeEmptyStr(data []string) (ret []string) {
	for _, v := range data {
		if v != "" {
			ret = append(ret, v)
		}
	}
	return ret
}

func copyAndCapture(r io.Reader) (string, error) {
	var laststring string
	buf := make([]byte, 1024, 1024)
	for {
		n, err := r.Read(buf[:])
		if n > 0 {
			d := buf[:n]
			laststring = string(d)
		}
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			items := removeEmptyStr(strings.Split(laststring, "\n"))
			if len(items) > 0 {
				return items[len(items)-1], nil
			}
			return "", errors.New("No data return")
		}
	}
}

func addMapToEnv(data map[string]interface{}) (newEnv []string) {
	newEnv = os.Environ()
	for v, k := range data {
		newEnv = append(newEnv, strings.ToUpper(v)+"="+k.(string))
	}
	return newEnv
}

// InterpretorTask ...
type InterpretorTask struct {
	Interpretor string   `json:"interpretor"`
	File        string   `json:"file"`
	Info        TaskInfo `json:"info"`
}

// Get ...
func (b *InterpretorTask) Get() ITask {
	return b
}

// GetID ...
func (b *InterpretorTask) GetID() string {
	return b.File
}

// GetInfo ...
func (b *InterpretorTask) GetInfo() TaskInfo {
	return b.Info
}

// Execute ...
func (b *InterpretorTask) Execute(file string, params map[string]interface{}, channel chan TaskFeedBack) {

	f := filepath.Join(getWorkingPath(), b.Interpretor, b.File)

	cmd := exec.Command(b.Interpretor, f, file)

	cmd.Env = addMapToEnv(params)

	output, _ := cmd.StdoutPipe()

	if err := cmd.Start(); err != nil {
		panic(err)
	}

	var wg sync.WaitGroup
	var lastLine string
	wg.Add(1)
	go func() {
		lastLine, _ = copyAndCapture(output)
		wg.Done()
	}()

	wg.Wait()
	err := cmd.Wait()
	if e, ok := err.(*exec.ExitError); ok {
		if !e.Success() {
		}
	}
	files := strings.Split(lastLine, ";")
	SendDone(channel, TaskOver{
		Files: removeEmptyStr(files),
	})
}
