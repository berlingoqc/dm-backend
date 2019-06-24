package program

import (
	"os/exec"
)

// RunningMode are the different way that a program can be run
type RunningMode string

const (
	// NoRunner is to not self run the program
	NoRunner RunningMode = "NO"
	// PathRunner to run to version of the program install on the computer
	PathRunner RunningMode = "PATH"
	// VersionRunner to run a version of the program download from a server to have a specif version
	// can retrieve from http server or github release
	VersionRunner RunningMode = "VERSION"
)

// Settings ...
type Settings struct {
	Mode RunningMode
}

// RunnerSettings ...
type RunnerSettings struct {
	Version string `json:"version"`
}

// Runner ...
type Runner struct {
	ErrorChan chan error
	Cmd       *exec.Cmd
}

// GetRunner ...
func GetRunner(program string, args map[string]interface{}, env map[string]interface{}) (*Runner, error) {
	ch := make(chan error)
	cmd := exec.Command(program)
	err := cmd.Start()
	if err != nil {
		return nil, err
	}
	go func() {
		err := cmd.Wait()
		if err != nil {
			print("Program end with error ", err.Error())
			ch <- err
			return
		}
		print("End without error")
	}()
	return &Runner{
		ErrorChan: ch,
		Cmd:       cmd,
	}, nil
}
