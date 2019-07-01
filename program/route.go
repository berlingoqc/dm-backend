package program

import "os/exec"

// RunnerInfo ...
type RunnerInfo struct {
	Name         string `json:"name"`
	Running      bool   `json:"state"`
	Error        string `json:"error"`
	OutputBuffer string `json:"outputbuffer"`
	ErrorBuffer  string `json:"errorbuffer"`

	Cmd      *exec.Cmd `json:"cmd"`
	Settings Settings  `json:"settings"`
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
