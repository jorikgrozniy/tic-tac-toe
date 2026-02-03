package middleware

import (
	"log"
	"net/http"
)

func Logger(mux *http.ServeMux) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)

		w.Header().Set("Content-Type", "application/json")

		mux.ServeHTTP(w, r)
	})
}
