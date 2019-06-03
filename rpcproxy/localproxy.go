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
			methodValue := reflect.ValueOf(handler).MethodByName(method[1])
			if methodValue.IsValid() {
				output := methodValue.Call(argsIn)
				println(output[0].String())

				response := make([]interface{}, len(output))
				for i, v := range output {
					response[i] = v.Interface()
				}
				rpcCall.Result = response
				return json.Marshal(rpcCall)
			}
			return nil, errors.New("method not found in namespace")
		}
		return nil, errors.New("namespace not found")
	}
	return nil, errors.New("no namespace")
}

// System ...
type System struct {
	Notification []string `json:"notification"`
	Methods      []string `json:"methods"`
	Namesapce    []string `json:"namespace"`
}

// ListNamespace ...
func (s *System) ListNamespace() []string {
	return s.Namesapce
}

// ListNotifications ...
func (s *System) ListNotifications() []string {
	return s.Notification
}

// ListMethods ...
func (s *System) ListMethods() []string {
	return s.Methods
}

// RegisterLocalHandler ...
func RegisterLocalHandler(ns string, lh *LocalHandler) {
	s := &System{
		Notification: make([]string, 0),
	}
	lh.Handlers["system"] = s
	for subns, obj := range lh.Handlers {
		fooType := reflect.TypeOf(obj)
		for i := 0; i < fooType.NumMethod(); i++ {
			method := fooType.Method(i)
			s.Methods = append(s.Methods, subns+"."+method.Name)
		}
		s.Namesapce = append(s.Namesapce, subns)
	}

	Handlers[ns] = lh
}
