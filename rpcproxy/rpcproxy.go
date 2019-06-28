package rpcproxy

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"

	"encoding/json"

	"github.com/gorilla/mux"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/cors"
)

// Handlers ...
var Handlers = make(map[string]RPCHandler)

// handleRPCProxy ...
func handleRPCProxy(w http.ResponseWriter, r *http.Request) {
	addCors(w)
	var err error
	defer func() {
		if err != nil {
			w.WriteHeader(200)
			w.Write(ErrorCall(err))
		}
	}()
	if da := r.Header.Get("X-dm-namespace"); da != "" {
		if handler, ok := Handlers[da]; ok {
			reader := r.Body
			var body []byte
			body, err = ioutil.ReadAll(reader)
			if err != nil {
				return
			}
			body, err = handler.Handle(body)
			if err != nil {
				return
			}

			w.WriteHeader(200)
			w.Write(body)
		}
	} else {
		println("PAS DE HEADER X-dm-namespace")
	}
}

// ProxyRPCRequest ...
func ProxyRPCRequest(host string, body []byte) ([]byte, error) {
	u := url.URL{Scheme: "http", Host: host, Path: "/jsonrpc"}

	resp, err := http.Post(u.String(), "appplication/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

// RPCRequest ...
func RPCRequest(host string, call RPCCall, result interface{}) error {
	u := url.URL{Scheme: "http", Host: host, Path: "/jsonrpc"}
	body, err := json.Marshal(call)
	if err != nil {
		return err
	}
	resp, err := http.Post(u.String(), "appplication/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	rpcCall := RPCCall{}
	err = json.Unmarshal(body, &rpcCall)
	if err != nil {
		return err
	}
	return mapstructure.Decode(rpcCall.Result[0], result)
}

// Register this
func Register(mux *mux.Router) http.Handler {
	mux.Path("/jsonrpc").Methods("POST").HandlerFunc(handleRPCProxy)
	mux.Path("/jsonrpc").Methods("GET").HandlerFunc(handleWebSocket)

	// Apres le registering enregistre son module local de system pour les infos
	// sur les namespace pour la reflection
	localHandler := &LocalHandler{
		Handlers: make(map[string]interface{}),
	}
	s := &System{
		Notification: make([]string, 0),
	}
	localHandler.Handlers["system"] = s
	Handlers["system"] = localHandler

	// Ajout la liste des methodes de system
	ns := "system"
	fooType := reflect.TypeOf(s)
	for i := 0; i < fooType.NumMethod(); i++ {
		method := fooType.Method(i)
		s.Methods = append(s.Methods, ns+"."+method.Name)
	}

	for k := range Handlers {
		s.Namesapce = append(s.Namesapce, k)
	}

	return cors.Default().Handler(mux)
}
