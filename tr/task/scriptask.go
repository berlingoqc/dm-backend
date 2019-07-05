package task

import (
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
)

func copyAndCapture(r io.Reader) (string, error) {
	var lastline string
	buf := make([]byte, 1024, 1024)
	for {
		n, err := r.Read(buf[:])
		if n > 0 {
			d := buf[:n]
			println(d)
		}
		if err != nil {
			// Read returns io.EOF at the end of file, which is not an error for us
			if err == io.EOF {
				err = nil
			}
			return lastline, err
		}
	}
}

// InterpretorTask ...
type InterpretorTask struct {
	Interpretor string
	File        string
	info        TaskInfo
}

// Get ...
func (b *InterpretorTask) Get() ITask {
	return &InterpretorTask{}
}

// GetID ...
func (b *InterpretorTask) GetID() string {
	return b.info.Name
}

// GetInfo ...
func (b *InterpretorTask) GetInfo() TaskInfo {
	return b.info
}

// Execute ...
func (b *InterpretorTask) Execute(file string, params map[string]interface{}, channel chan TaskFeedBack) {
	cmd := exec.Command(b.Interpretor, b.File, file)

	//cmd.Env = addMapToEnv(params)

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
		Files: files,
	})
}

func addMapToEnv(data map[string]interface{}) (newEnv []string) {
	newEnv = os.Environ()
	return newEnv
}
