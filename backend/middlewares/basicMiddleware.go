package middlewares

import (
	"log"
	"net/http"
)

func BasicMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(wr http.ResponseWriter, req *http.Request) {
		log.Println("Middleware called on", req.URL.Path)
		// do stuff
		h.ServeHTTP(wr, req)
	})
}
