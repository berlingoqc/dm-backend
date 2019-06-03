package aria2

import (
	"github.com/berlingoqc/dm/rpcproxy"
	"github.com/berlingoqc/dm/tr"
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

func init() {
	rpcproxy.Handlers["aria2"] = &RPCHandler{}
	tr.Handlers["aria2"] = &FileHandler{}
}
