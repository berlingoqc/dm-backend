package pipeline

// IFileHandler ...
type IFileHandler interface {
	GetFilePath(i []interface{}) (string, error)
	GetEvent(method string) string
	SetConfig(data interface{}) 
}

// Handlers ...
var Handlers = make(map[string]IFileHandler)
