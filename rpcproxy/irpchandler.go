package rpcproxy

// RPCCall ...
type RPCCall struct {
	Jsonrpc string        `json:"jsonrpc"`
	ID      string        `json:"id"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	Result  []interface{} `json:"result"`
	Error   interface{}   `json:"error"`
}

// RPCHandlerEndpoint ...
type RPCHandlerEndpoint struct {
	Namespace string
	URL       string
}

// RPCHandler ...
type RPCHandler interface {
	GetConfig() *RPCHandlerEndpoint
	SetConfig(*RPCHandlerEndpoint)
	Handle(body []byte) ([]byte, error)
}
