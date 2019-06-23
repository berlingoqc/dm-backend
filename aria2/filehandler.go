package aria2

import (
	"errors"

	"github.com/berlingoqc/dm-backend/rpcproxy"
)

// FileInfo ...
type FileInfo struct {
	CompletedLength string `json:"completedLength"`
	Path            string `json:"path"`
}

// DownloadOver ...
type DownloadOver struct {
	Gid string `json:"gid"`
}

// FileHandler ...
type FileHandler struct {
	config *rpcproxy.RPCHandlerEndpoint
}

// GetEvent ...
func (f *FileHandler) GetEvent(method string) string {
	switch method {
	case "aria2.onDownloadComplete":
		return "onDownloadOver"
	}
	return ""
}

// GetFilePath ...
func (f *FileHandler) GetFilePath(i []interface{}) (string, error) {
	if len(i) != 1 {
		return "", errors.New("CACA")
	}
	gid := i[0].(map[string]interface{})["gid"].(string)
	println("GID ", gid)

	// Doit faire une request avec le id pour avoir le full path du fichier
	fileInfo := FileInfo{}
	if err := rpcproxy.RPCRequest(f.config.URL, rpcproxy.RPCCall{
		Jsonrpc: "2.0",
		ID:      "qwer",
		Method:  "aria2.getFiles",
		Params:  []interface{}{gid},
	}, &fileInfo); err != nil {
		return "", err
	}
	return fileInfo.Path, nil
}

// GetDownloadOverMessage ...
func (f *FileHandler) GetDownloadOverMessage() interface{} {
	return &DownloadOver{}
}

// SetConfig ...
func (f *FileHandler) SetConfig(data interface{}) {
	f.config = data.(*rpcproxy.RPCHandlerEndpoint)
}
