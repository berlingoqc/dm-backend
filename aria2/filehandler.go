package aria2

import (
	"github.com/berlingoqc/dm-backend/rpcproxy"
)

// FileInfo ...
type FileInfo struct {
	CompletedLength string `json:"completedLength"`
	Path            string `json:"path"`
}

// FileHandler ...
type FileHandler struct {
	Config *rpcproxy.RPCHandlerEndpoint
}

// GetFile ...
func (f *FileHandler) GetFile(d interface{}, expected interface{}) (string, bool, error) {
	dd := d.([]interface{})
	gid := dd[0].(map[string]interface{})["gid"].(string)
	println("GID ", gid)
	if gid != expected.(string) {
		return "", false, nil
	}
	// Doit faire une request avec le id pour avoir le full path du fichier
	fileInfo := FileInfo{}
	if err := rpcproxy.RPCRequest(f.Config.URL, rpcproxy.RPCCall{
		Jsonrpc: "2.0",
		ID:      "qwer",
		Method:  "aria2.getFiles",
		Params:  []interface{}{gid},
	}, &fileInfo); err != nil {
		return "", false, err
	} else {
		return fileInfo.Path, true, err
	}
}
