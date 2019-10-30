package aria2

import (
	"github.com/berlingoqc/dm-backend/rpcproxy"
)

// RPCHandler ...
type RPCHandler struct {
	config *rpcproxy.RPCHandlerEndpoint
}

// GetConfig ...
func (a *RPCHandler) GetConfig() *rpcproxy.RPCHandlerEndpoint {
	return a.config
}

// SetConfig ...
func (a *RPCHandler) SetConfig(config *rpcproxy.RPCHandlerEndpoint) {
	a.config = config
}

// Handle ...
func (a *RPCHandler) Handle(body []byte) ([]byte, error) {
	return rpcproxy.ProxyRPCRequest(a.config.URL, body)
}

func init() {
	rpcproxy.Handlers["aria2"] = &RPCHandler{}
}
