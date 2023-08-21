package middleware

import (
	"net/http"
)

// Обворачивает функцию next, вызывая ее только, если req.Method == method
func WithMethod(method string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.Method == method {
			next.ServeHTTP(w, req)
		}
	})
}
