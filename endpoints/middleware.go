package endpoints

import (
	"net/http"
)

// HandleCORS is a CORS handler.
func HandleCORS(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")
	if origin == "" {
		origin = "*"
	}

	w.Header().Set("Access-Control-Allow-Origin", origin)
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, PATCH")
	w.Header().Set("Access-Control-Allow-Headers",
		"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	// Stop if the request is OPTIONS.
	if r.Method == "OPTIONS" {
		return
	}
}

// CORSMiddleware sets up CORS headers.
func CORSMiddleware(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		HandleCORS(w, r)
		if r.Method == "OPTIONS" {
			return
		}
		// Call the next handler.
		handler.ServeHTTP(w, r)
	}
}
