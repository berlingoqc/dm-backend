package rpcproxy

import (
	"encoding/json"
	"errors"
	"reflect"
	"strings"
)

// LocalHandler ...
type LocalHandler struct {
	Handlers map[string]interface{}
}

// GetConfig ...
func (a *LocalHandler) GetConfig() *RPCHandlerEndpoint {
	return nil
}

// SetConfig ...
func (a *LocalHandler) SetConfig(config *RPCHandlerEndpoint) {}

// Handle ...
func (a *LocalHandler) Handle(body []byte) ([]byte, error) {
	rpcCall := RPCCall{}
	if err := json.Unmarshal(body, &rpcCall); err != nil {
		return nil, err
	}
	method := strings.Split(rpcCall.Method, ".")
	if len(method) > 1 {
		if handler, ok := a.Handlers[method[0]]; ok {
			argsIn := make([]reflect.Value, len(rpcCall.Params))
			for i, v := range rpcCall.Params {
				argsIn[i] = reflect.ValueOf(v)
			}
			output := reflect.ValueOf(handler).MethodByName(method[1]).Call(argsIn)
			println(output[0].String())

			response := make([]interface{}, len(output))
			for i, v := range output {
				response[i] = v.Interface()
			}
			rpcCall.Result = response
			return json.Marshal(rpcCall)
		}
		return nil, errors.New("namespace not found")
	}
	return nil, errors.New("no namespace")
}

type system struct {
	Notification []string
	Methods      []string
}

func (s system) listNotifications() []string {
	return s.Notification
}

func (s system) listMethods() []string {
	return s.Methods
}
