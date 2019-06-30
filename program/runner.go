package program

import (
	"errors"
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
	Name string                 `json:"name"`
	Mode RunningMode            `json:"mode"`
	Args []string               `json:"args"`
	Env  map[string]interface{} `json:"env"`
}

// RunnerSettings ...
type RunnerSettings struct {
	Version string `json:"version"`
}

// Runner ...
type Runner struct {
	ErrorChan chan error
	Cmd       *exec.Cmd
	Running   bool
}

// getRunner ...
func getRunner(program string, args []string, env map[string]interface{}) (*Runner, error) {
	ch := make(chan error)
	cmd := exec.Command(program, args...)
	err := cmd.Start()
	if err != nil {
		return nil, err
	}
	r := &Runner{
		ErrorChan: ch,
		Cmd:       cmd,
		Running:   true,
	}
	go func() {
		err := cmd.Wait()
		r.Running = false
		if err != nil {
			ch <- err
			return
		}
		print("End without error")
	}()
	return r, nil
}

func programListen() {
	for {
		for _, v := range activeRunner {
			select {
			case err := <-v.ErrorChan:
				println("RUNNER ERROR ", err.Error())
				break
			}

		}
	}

}

// GetRunner ...
func GetRunner(s *Settings) (*Runner, error) {
	switch s.Mode {
	case PathRunner:
		return getRunner(s.Name, s.Args, s.Env)
	default:
		return nil, errors.New("Mode not impement yet")
	}
}

var activeRunner = make(map[string]*Runner)

// Start ...
func Start(program []*Settings) error {
	for _, v := range program {
		r, err := GetRunner(v)
		if err != nil {
			return err
		}
		activeRunner[v.Name] = r
	}
	return nil
}

func init() {
	go programListen()
}
