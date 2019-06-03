package rpcproxy

import "encoding/json"

// ErrorCall ...
func ErrorCall(err error) []byte {
	call := RPCCall{
		Jsonrpc: "2.0",
		ID:      "qwer",
		Error:   err.Error(),
	}

	d, _ := json.Marshal(&call)
	return d
}
