package ydl

import "github.com/berlingoqc/dm-backend/rpcproxy"

// FileHandler ...
type FileHandler struct {
	config *rpcproxy.RPCHandlerEndpoint
}

// DownloadOver ...
type DownloadOver struct{}

// GetEvent ...
func (f *FileHandler) GetEvent(method string) string {
	switch method {
	case "onDownloadOver":
		return ""
	}
	return ""
}

// GetFilePath ...
func (f *FileHandler) GetFilePath(i []interface{}) (string, error) {
	return "", nil
}

// GetDownloadOverMessage ...
func (f *FileHandler) GetDownloadOverMessage() interface{} {
	return &DownloadOver{}
}

// SetConfig ...
func (f *FileHandler) SetConfig(data interface{}) {
	f.config = data.(*rpcproxy.RPCHandlerEndpoint)
}
