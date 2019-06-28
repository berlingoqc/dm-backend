package rpcproxy

import "net/http"

func addCors(w http.ResponseWriter) {

	//Allow CORS here By * or specific origin
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
	w.Header().Set("Access-Control-Request-Method", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, X-dm-namespace")
}
