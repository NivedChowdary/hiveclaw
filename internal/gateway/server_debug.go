package gateway

import (
	"log"
	"net/http"
)

// DebugHandler wraps file server with logging
func DebugFileServer(fs http.FileSystem) http.Handler {
	fileServer := http.FileServer(fs)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[DEBUG] Request: %s %s", r.Method, r.URL.Path)
		fileServer.ServeHTTP(w, r)
	})
}
