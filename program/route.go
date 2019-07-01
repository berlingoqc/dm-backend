package program

import "os/exec"

// RunnerInfo ...
type RunnerInfo struct {
	Name         string `json:"name"`
	Running      bool   `json:"state"`
	Error        string `json:"error"`
	OutputBuffer string `json:"outputbuffer"`
	ErrorBuffer  string `json:"errorbuffer"`

	Cmd *exec.Cmd `json:"cmd"`
}

// RPC ...
type RPC struct{}

// GetActiveRunner ...
func (r *RPC) GetActiveRunner() []RunnerInfo {
	var infos []RunnerInfo
	for k, v := range activeRunner {
		r := RunnerInfo{
			Name:         k,
			Running:      v.Running,
			OutputBuffer: string(v.STDOut.Bytes()),
			ErrorBuffer:  string(v.STDErr.Bytes()),
			Cmd:          v.Cmd,
		}
		if v.Error != nil {
			r.Error = v.Error.Error()
		}

		infos = append(infos, r)
	}
	return infos
}

// StartRunner ...
func (r *RPC) StartRunner(name string) string {
	if set, ok := settingsRunner[name]; ok {
		delete(activeRunner, name)
		if err := Start([]*Settings{set}); err != nil {
			panic(err)
		}
	}
	return "OK"
}

// StopRunner ...
func (r *RPC) StopRunner(name string) string {
	if r, ok := activeRunner[name]; ok {
		if err := r.Cmd.Process.Kill(); err != nil {
			panic(err)
		}
	}
	return "OK"
}
