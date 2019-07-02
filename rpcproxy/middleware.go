package rpcproxy

import "net/http"

func addCors(w http.ResponseWriter) {
	//Allow CORS here By * or specific origin
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
	w.Header().Set("Access-Control-Request-Method", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, X-dm-namespace")
}

// SecurityHeader ...
var SecurityHeader = "X-dm-auth"

// ValidToken ...
var ValidToken func(token string, r *http.Request) error

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get(SecurityHeader)
		if err := ValidToken(token, r); err != nil {
			http.Error(w, "Forbidden", http.StatusForbidden)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}
