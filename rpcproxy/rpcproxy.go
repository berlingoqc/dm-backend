package rpcproxy

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"

	"encoding/json"

	"github.com/gorilla/mux"
	"github.com/mitchellh/mapstructure"
)

// Handlers ...
var Handlers = make(map[string]RPCHandler)

// handleRPCProxy ...
func handleRPCProxy(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		if err != nil {
			w.WriteHeader(400)
			w.Write([]byte("ERROR"))
		}
	}()
	if da := r.Header.Get("X-dm-namespace"); da != "" {
		if handler, ok := Handlers[da]; ok {
			println("Proxing for ", da)
			reader := r.Body
			body, err := ioutil.ReadAll(reader)

			u := url.URL{Scheme: "http", Host: handler.GetConfig().URL, Path: "/jsonrpc"}

			resp, err := http.Post(u.String(), "appplication/json", bytes.NewBuffer(body))
			if err != nil {
				println(err)
				return
			}
			defer resp.Body.Close()
			body, err = ioutil.ReadAll(resp.Body)
			if err != nil {
				println(err)
				return
			}
			w.WriteHeader(resp.StatusCode)
			w.Write(body)
		}
	} else {
		println("PAS DE HEADER X-dm-namespace")
	}
}

// RPCRequest ...
func RPCRequest(host string, call RPCCall, result interface{}) error {
	u := url.URL{Scheme: "http", Host: host, Path: "/jsonrpc"}
	body, err := json.Marshal(call)
	if err != nil {
		return err
	}
	println(string(body))
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
func Register(mux *mux.Router) {
	mux.Path("/jsonrpc").Methods("POST").HandlerFunc(handleRPCProxy)
	mux.Path("/jsonrpc").Methods("GET").HandlerFunc(handleWebSocket)
}
